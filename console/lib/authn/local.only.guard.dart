import 'package:flutter/material.dart';
import 'package:retrovibed/meta.dart' as _meta;
import 'package:retrovibed/authz.dart' as authz;
import 'cache.dart';

// Decides whether guarded should even be mounted: a local_only
// (guest) identity should never make the deeppool calls that
// widget performs, so it's simplest to never create its state at all rather
// than guard each call inside it.
//
// Listens to AuthzCache's `changed` notifier rather than depending on
// AuthzTokenData (InheritedWidget) directly - refresh() only ever reassigns
// meta.current in place, not the AuthzTokenData.meta reference itself, so
// InheritedWidget dependents aren't notified past the first resolved token.
// `changed.value = bearer` fires correctly on every refresh, so it's the
// reliable signal to rebuild on here.
class LocalOnlyGuard extends StatelessWidget {
  final Widget child;
  final Widget Function(Widget) guarded;
  const LocalOnlyGuard(this.child, this.guarded, {super.key});

  @override
  Widget build(BuildContext context) {
    final cache = AuthzCache.of(context);
    return ValueListenableBuilder<authz.Bearer<_meta.Token>>(
      valueListenable: cache.changed,
      builder: (context, bearer, _) {
        if (bearer.token.localOnly) {
          return child;
        }
        return guarded(child);
      },
    );
  }
}
