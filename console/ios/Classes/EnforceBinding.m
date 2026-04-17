// Prevents the linker from stripping Go-exported symbols as dead code.
// These functions are called via DynamicLibrary.process() at runtime.
#import "libretrovibed.h"

__attribute__((used)) void enforce_binding() {
	oauth2_bearer();
	authn_bearer();
	authn_bearer_host("");
	public_key();
	username();
	ips();
	egdaemon("[]");
}
