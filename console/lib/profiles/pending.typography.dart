import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/meta.dart' as meta;
import 'cache.dart';

class Typography extends StatelessWidget {
  final meta.Profile current;
  final List<Widget> leading;
  final List<Widget> trailing;

  const Typography(
    this.current, {
    super.key,
    this.leading = const [],
    this.trailing = const [],
  });

  static Widget removebtn(BuildContext context, String id, {VoidCallback? onPressed}) {
    return ds.buttons.remove(
      onPressed: onPressed,
    );
  }

  static Widget approvebtn(BuildContext context, String id, {VoidCallback? onPressed}) {
    return ds.buttons.accept(
      onPressed: onPressed,
    );
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
    return Row(
      spacing: defaults.spacing,
      children: [
        ...leading,
        Expanded(child: Text(current.id)),
        Expanded(child: Text(current.display.isEmpty ? "-" : current.display)),
        Expanded(child: ds.Timestamp.iso8601(current.updatedAt)),
        ...trailing,
      ],
    );
  }
}
