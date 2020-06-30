package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/watchtower/api"
	"github.com/hashicorp/watchtower/globals"
	"github.com/hashicorp/watchtower/internal/gen/controller/api/services"
	"github.com/hashicorp/watchtower/internal/servers/controller/handlers"
	"github.com/hashicorp/watchtower/internal/servers/controller/handlers/groups"
	"github.com/hashicorp/watchtower/internal/servers/controller/handlers/host_catalogs"
	"github.com/hashicorp/watchtower/internal/servers/controller/handlers/host_sets"
	"github.com/hashicorp/watchtower/internal/servers/controller/handlers/hosts"
	"github.com/hashicorp/watchtower/internal/servers/controller/handlers/organizations"
	"github.com/hashicorp/watchtower/internal/servers/controller/handlers/projects"
	"github.com/hashicorp/watchtower/internal/servers/controller/handlers/roles"
	"github.com/hashicorp/watchtower/internal/servers/controller/handlers/users"
	"github.com/hashicorp/watchtower/internal/ui"
)

type HandlerProperties struct {
	ListenerConfig *configutil.Listener
}

// Handler returns an http.Handler for the services. This can be used on
// its own to mount the Vault API within another web server.
func (c *Controller) handler(props HandlerProperties) (http.Handler, error) {
	// Create the muxer to handle the actual endpoints
	mux := http.NewServeMux()

	h, err := handleGrpcGateway(c)
	if err != nil {
		return nil, err
	}
	mux.Handle("/v1/", h)

	mux.Handle("/", handleUi(c))

	corsWrappedHandler := wrapHandlerWithCors(mux, props)
	commonWrappedHandler := wrapHandlerWithCommonFuncs(corsWrappedHandler, c, props)

	return commonWrappedHandler, nil
}

func handleUi(c *Controller) http.Handler {
	var nextHandler http.Handler
	if c.conf.RawConfig.PassthroughDirectory != "" {
		nextHandler = ui.DevPassthroughHandler(c.logger, c.conf.RawConfig.PassthroughDirectory)
	} else {
		nextHandler = http.FileServer(ui.AssetFile())
	}

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			irw := newIndexResponseWriter(c.conf.DefaultOrgId)
			nextHandler.ServeHTTP(irw, r)
			irw.writeToWriter(w)

		default:
			nextHandler.ServeHTTP(w, r)
		}
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		dotIndex := strings.LastIndex(r.URL.Path, ".")
		switch dotIndex {
		case -1:
			// For all paths without an extension serve /index.html
			r.URL.Path = "/"

		default:
			switch r.URL.Path {
			case "/", "/favicon.png", "/assets/styles.css":

			default:
				for i := dotIndex + 1; i < len(r.URL.Path); i++ {
					intVal := r.URL.Path[i]
					// Current guidance from FE is if it's only alphanum after
					// the last dot, treat it as an extension
					if intVal < '0' ||
						(intVal > '9' && intVal < 'A') ||
						(intVal > 'Z' && intVal < 'a') ||
						intVal > 'z' {
						// Not an extension. Serve the contents of index.html
						r.URL.Path = "/"
					}
				}
			}
		}

		// Fall through to the next handler
		rootHandler.ServeHTTP(w, r)
	})
}

func handleGrpcGateway(c *Controller) (http.Handler, error) {
	// Register*ServiceHandlerServer methods ignore the passed in ctx.  Using the baseContext now just in case this changes
	// in the future, at which point we'll want to be using the baseContext.
	ctx := c.baseContext
	mux := runtime.NewServeMux(runtime.WithProtoErrorHandler(handlers.ErrorHandler(c.logger)))
	hcs, err := host_catalogs.NewService(c.StaticHostRepoFn)
	if err != nil {
		return nil, fmt.Errorf("failed to create host catalog handler service: %w", err)
	}
	if err := services.RegisterHostCatalogServiceHandlerServer(ctx, mux, hcs); err != nil {
		return nil, fmt.Errorf("failed to register host catalog service handler: %w", err)
	}
	if err := services.RegisterHostSetServiceHandlerServer(ctx, mux, &host_sets.Service{}); err != nil {
		return nil, fmt.Errorf("failed to register host set service handler: %w", err)
	}
	if err := services.RegisterHostServiceHandlerServer(ctx, mux, &hosts.Service{}); err != nil {
		return nil, fmt.Errorf("failed to register host service handler: %w", err)
	}
	os, err := organizations.NewService(c.IamRepoFn)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization handler service: %w", err)
	}
	if err := services.RegisterOrganizationServiceHandlerServer(ctx, mux, os); err != nil {
		return nil, fmt.Errorf("failed to register organization service handler: %w", err)
	}
	ps, err := projects.NewService(c.IamRepoFn)
	if err != nil {
		return nil, fmt.Errorf("failed to create project handler service: %w", err)
	}
	if err := services.RegisterProjectServiceHandlerServer(ctx, mux, ps); err != nil {
		return nil, fmt.Errorf("failed to register project service handler: %w", err)
	}
	us, err := users.NewService(c.IamRepoFn)
	if err != nil {
		return nil, fmt.Errorf("failed to create user handler service: %w", err)
	}
	if err := services.RegisterUserServiceHandlerServer(ctx, mux, us); err != nil {
		return nil, fmt.Errorf("failed to register user service handler: %w", err)
	}
	gs, err := groups.NewService(c.IamRepoFn)
	if err != nil {
		return nil, fmt.Errorf("failed to create group handler service: %w", err)
	}
	if err := services.RegisterGroupServiceHandlerServer(ctx, mux, gs); err != nil {
		return nil, fmt.Errorf("failed to register group service handler: %w", err)
	}
	rs, err := roles.NewService(c.IamRepoFn)
	if err != nil {
		return nil, fmt.Errorf("failed to create role handler service: %w", err)
	}
	if err := services.RegisterRoleServiceHandlerServer(ctx, mux, rs); err != nil {
		return nil, fmt.Errorf("failed to register role service handler: %w", err)
	}

	return mux, nil
}

func wrapHandlerWithCommonFuncs(h http.Handler, c *Controller, props HandlerProperties) http.Handler {
	var maxRequestDuration time.Duration
	var maxRequestSize int64
	if props.ListenerConfig != nil {
		maxRequestDuration = props.ListenerConfig.MaxRequestDuration
		maxRequestSize = props.ListenerConfig.MaxRequestSize
	}
	if maxRequestDuration == 0 {
		maxRequestDuration = globals.DefaultMaxRequestDuration
	}
	if maxRequestSize == 0 {
		maxRequestSize = globals.DefaultMaxRequestSize
	}
	var defaultOrgId string
	if c != nil && c.conf != nil {
		defaultOrgId = c.conf.DefaultOrgId
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if defaultOrgId != "" {
			splitPath := strings.Split(r.URL.Path, "/")
			if len(splitPath) >= 3 && splitPath[2] == "projects" {
				http.Redirect(w, r, path.Join("/v1/orgs", defaultOrgId, strings.Join(splitPath[2:], "/")), 307)
				return
			}
		}

		// Set the Cache-Control header for all responses returned
		w.Header().Set("Cache-Control", "no-store")

		// Start with the request context
		ctx := r.Context()
		var cancelFunc context.CancelFunc
		// Add our timeout
		ctx, cancelFunc = context.WithTimeout(ctx, maxRequestDuration)
		// Add a size limiter if desired
		if maxRequestSize > 0 {
			ctx = context.WithValue(ctx, "max_request_size", maxRequestSize)
		}
		ctx = context.WithValue(ctx, "original_request_path", r.URL.Path)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
		cancelFunc()
		return
	})
}

func wrapHandlerWithCors(h http.Handler, props HandlerProperties) http.Handler {
	allowedMethods := []string{
		http.MethodDelete,
		http.MethodGet,
		http.MethodOptions,
		http.MethodPost,
		http.MethodPatch,
	}

	allowedOrigins := props.ListenerConfig.CorsAllowedOrigins

	allowedHeaders := append([]string{
		"Content-Type",
		"X-Requested-With",
		"Authorization",
	}, props.ListenerConfig.CorsAllowedHeaders...)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !props.ListenerConfig.CorsEnabled {
			h.ServeHTTP(w, req)
			return
		}

		origin := req.Header.Get("Origin")

		if origin == "" {
			// Serve directly
			h.ServeHTTP(w, req)
			return
		}

		// Check origin
		var valid bool
		switch {
		case len(allowedOrigins) == 0:
			// not valid

		case len(allowedOrigins) == 1 && allowedOrigins[0] == "*":
			valid = true

		default:
			valid = strutil.StrListContains(allowedOrigins, origin)
		}
		if !valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)

			err := &api.Error{
				Status: api.Int(http.StatusForbidden),
				Code:   api.String("origin forbidden"),
			}

			enc := json.NewEncoder(w)
			enc.Encode(err)
			return
		}

		if req.Method == http.MethodOptions &&
			!strutil.StrListContains(allowedMethods, req.Header.Get("Access-Control-Request-Method")) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")

		// Apply headers for preflight requests
		if req.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
			w.Header().Set("Access-Control-Max-Age", "300")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, req)
		return
	})
}

/*
func WrapForwardedForHandler(h http.Handler, authorizedAddrs []*sockaddr.SockAddrMarshaler, rejectNotPresent, rejectNonAuthz bool, hopSkips int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers, headersOK := r.Header[textproto.CanonicalMIMEHeaderKey("X-Forwarded-For")]
		if !headersOK || len(headers) == 0 {
			if !rejectNotPresent {
				h.ServeHTTP(w, r)
				return
			}
			respondError(w, http.StatusBadRequest, fmt.Errorf("missing x-forwarded-for header and configured to reject when not present"))
			return
		}

		host, port, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// If not rejecting treat it like we just don't have a valid
			// header because we can't do a comparison against an address we
			// can't understand
			if !rejectNotPresent {
				h.ServeHTTP(w, r)
				return
			}
			respondError(w, http.StatusBadRequest, errwrap.Wrapf("error parsing client hostport: {{err}}", err))
			return
		}

		addr, err := sockaddr.NewIPAddr(host)
		if err != nil {
			// We treat this the same as the case above
			if !rejectNotPresent {
				h.ServeHTTP(w, r)
				return
			}
			respondError(w, http.StatusBadRequest, errwrap.Wrapf("error parsing client address: {{err}}", err))
			return
		}

		var found bool
		for _, authz := range authorizedAddrs {
			if authz.Contains(addr) {
				found = true
				break
			}
		}
		if !found {
			// If we didn't find it and aren't configured to reject, simply
			// don't trust it
			if !rejectNonAuthz {
				h.ServeHTTP(w, r)
				return
			}
			respondError(w, http.StatusBadRequest, fmt.Errorf("client address not authorized for x-forwarded-for and configured to reject connection"))
			return
		}

		// At this point we have at least one value and it's authorized

		// Split comma separated ones, which are common. This brings it in line
		// to the multiple-header case.
		var acc []string
		for _, header := range headers {
			vals := strings.Split(header, ",")
			for _, v := range vals {
				acc = append(acc, strings.TrimSpace(v))
			}
		}

		indexToUse := len(acc) - 1 - hopSkips
		if indexToUse < 0 {
			// This is likely an error in either configuration or other
			// infrastructure. We could either deny the request, or we
			// could simply not trust the value. Denying the request is
			// "safer" since if this logic is configured at all there may
			// be an assumption it can always be trusted. Given that we can
			// deny accepting the request at all if it's not from an
			// authorized address, if we're at this point the address is
			// authorized (or we've turned off explicit rejection) and we
			// should assume that what comes in should be properly
			// formatted.
			respondError(w, http.StatusBadRequest, fmt.Errorf("malformed x-forwarded-for configuration or request, hops to skip (%d) would skip before earliest chain link (chain length %d)", hopSkips, len(headers)))
			return
		}

		r.RemoteAddr = net.JoinHostPort(acc[indexToUse], port)
		h.ServeHTTP(w, r)
		return
	})
}
*/
