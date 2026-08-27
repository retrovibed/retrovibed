import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import './api.dart' as api;

class KnownMediaTypography extends StatelessWidget {
  final api.Known current;
  final GestureTapCallback? onDoubleTap;
  final List<Widget> leading;
  final List<Widget> trailing;
  final BoxDecoration? decoration;

  const KnownMediaTypography(
    this.current, {
    super.key,
    this.onDoubleTap,
    this.leading = const [],
    this.trailing = const [],
    this.decoration,
  });

  static Widget removebtn(BuildContext context, String id, {VoidCallback? onPressed}) {
    return ds.buttons.remove(
      onPressed: onPressed,
    );
  }

  static Widget fromID(
    String id, {
    List<Widget> leading = const [],
    List<Widget> trailing = const [],
    BoxDecoration? decoration,
  }) {
    return Builder(
      builder: (context) {
        return FutureBuilder<api.Known>(
          initialData: api.Known.create(),
          future: api.known
              .cached(
                id,
                () => api.known.get(
                  id,
                  options: [authn.request(authn.AuthzCache.meta(context))],
                ),
              )
              .then((w) => w.known),
          builder: (BuildContext ctx, AsyncSnapshot<api.Known> snapshot) {
            return ds.Loading(
              loading: !(snapshot.hasData || snapshot.hasError),
              cause: ds.Error.maybeErr(snapshot.error),
              snapshot.data == null
                  ? SizedBox()
                  : KnownMediaTypography(
                    snapshot.data!,
                    leading: leading,
                    trailing: trailing,
                    decoration: decoration,
                  ),
            );
          },
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);
    return ds.Container(
      padding: defaults.padding,
      decoration: decoration ?? const BoxDecoration(),
      Row(
        spacing: defaults.spacing,
        children: [
          ...leading,
          ds.Image(current.image, size: theme.textTheme.bodyMedium?.fontSize),
          Text(current.description),
          ...trailing,
        ],
      ),
    );
  }
}
