import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
// import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/timex.dart' as timex;
import 'cache.dart';
import 'disabled.icon.dart';
// import 'permissions.editor.dart';
// import 'rename.modal.dart';

class Typography extends StatelessWidget {
  static void noop(Future<meta.Profile> pending) {}
  final meta.Profile current;
  final void Function(Future<meta.Profile> pending) onChange;
  final GestureTapCallback? onDoubleTap;
  final List<Widget> leading;
  final List<Widget> trailing;

  const Typography(
    this.current, {
    super.key,
    this.onChange = noop,
    this.onDoubleTap,
    this.leading = const [],
    this.trailing = const [],
  });

  static Widget removebtn(
    BuildContext context,
    String id, {
    VoidCallback? onPressed,
  }) {
    return ds.buttons.remove(onPressed: onPressed);
  }

  static Widget fromID(
    String id, {
    List<Widget> leading = const [],
    List<Widget> trailing = const [],
  }) {
    return Builder(
      builder: (context) {
        return FutureBuilder<meta.Profile>(
          initialData: meta.Profile.create(),
          future: cached(
            id,
            () => meta.profiles.find(
              id,
              options: [authn.Authenticated.bearer(context)],
            ),
          ).then((w) => w.profile),
          builder: (BuildContext ctx, AsyncSnapshot<meta.Profile> snapshot) {
            return ds.Loading(
              loading: !(snapshot.hasData || snapshot.hasError),
              cause: ds.Error.maybeErr(snapshot.error),
              snapshot.data == null
                  ? SizedBox()
                  : Typography(
                    snapshot.data!,
                    leading: leading,
                    trailing: trailing,
                  ),
            );
          },
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final session = authn.Authenticated.syncSession(context);
    final defaultDisplay = current.display.isEmpty ? "-" : current.display;
    return Builder(
      builder: (context) {
        final compact = defaults.isCompact;
        return Material(
          color: Colors.transparent,
          child: Row(
            spacing: defaults.spacing,
            children: [
              ...leading,
              DisabledIcon(
                dates: [
                  timex.iso8601(current.disabledAt),
                  timex.iso8601(current.disabledManuallyAt),
                  timex.iso8601(current.disabledPendingApprovalAt),
                ],
                onTap: (enabled) {
                  final disable =
                      () => meta.profiles
                          .disable(
                            current.id,
                            options: [authn.AuthzCache.bearer(context)],
                          )
                          .then((v) => v.profile);
                  final enable =
                      () => meta.profiles
                          .enable(
                            current,
                            options: [authn.AuthzCache.bearer(context)],
                          )
                          .then((v) => v.profile);
                  final pending = enabled ? disable() : enable();
                  this.onChange(pending);
                  return pending;
                },
              ),
              if (!compact)
                Expanded(
                  child: Text(current.id, overflow: TextOverflow.ellipsis),
                ),
              Expanded(
                child: Text(
                  session.profile.id == current.id ? "you" : defaultDisplay,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              if (!compact)
                Expanded(child: ds.Timestamp.iso8601(current.updatedAt)),
              ...trailing,
            ],
          ),
        );
      },
    );
  }
}
