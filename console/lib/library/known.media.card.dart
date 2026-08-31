import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media/media.pb.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;
import './api.dart' as api;

class KnownMediaCard extends StatelessWidget {
  final api.Known current;
  final GestureTapCallback? onTap;
  final GestureTapCallback? onDoubleTap;
  final GestureTapCallback? onSecondaryTap;
  final GestureLongPressCallback? onLongPress;
  final ValueNotifier<bool>? hovered;
  final Widget help;
  final Widget? overlay;
  final List<Widget> leading;
  final List<Widget> trailing;
  final IconData? icon;
  final bool highlighted;
  final BoxConstraints? constraints;

  static const hint = ds.Hint(
    Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text("hover — reveal the media summary and primary action icon"),
        Text("tap — perform the primary action (e.g. play or download)"),
        Text("double tap — perform the secondary action"),
        Text("secondary tap — open the context menu"),
        Text("long press — show or hide additional controls"),
      ],
    ),
  );

  static Widget future(
    Future<api.Known> future, {
    Key? key,
    GestureTapCallback? onTap,
    GestureTapCallback? onDoubleTap,
    GestureTapCallback? onSecondaryTap,
    GestureLongPressCallback? onLongPress,
    List<Widget> leading = const [],
    List<Widget> trailing = const [],
    bool highlighted = false,
    IconData? icon = Icons.play_circle_filled,
    Widget help = ds.HelpScope.None,
    Widget? overlay,
    BoxConstraints? constraints,
  }) {
    return FutureBuilder<api.Known>(
      future: future,
      builder: (context, snapshot) {
        return KnownMediaCard(
          snapshot.data ?? api.Known(),
          key: key,
          onTap: onTap,
          onDoubleTap: onDoubleTap,
          onSecondaryTap: onSecondaryTap,
          onLongPress: onLongPress,
          leading: leading,
          trailing: trailing,
          highlighted: highlighted,
          icon: icon,
          help: help,
          overlay: overlay,
          constraints: constraints,
        );
      },
    );
  }

  static Widget auto(
    BuildContext context,
    Media m, {
    Key? key,
    GestureTapCallback? onTap,
    GestureTapCallback? onDoubleTap,
    GestureTapCallback? onSecondaryTap,
    GestureLongPressCallback? onLongPress,
    List<Widget> leading = const [],
    List<Widget> trailing = const [],
    bool highlighted = false,
    IconData? icon = Icons.play_circle_filled,
    Widget help = ds.HelpScope.None,
    Widget? overlay,
    BoxConstraints? constraints,
  }) {
    final authz = authn.AuthzCache.meta(context);
    return KnownMediaCard.future(
      api.known.autodetect(m, options: [authn.request(authz)]),
      key: key ?? ValueKey(uuidx.md5x("${m.id}.${m.updatedAt}")),
      onTap: onTap,
      onDoubleTap: onDoubleTap,
      onSecondaryTap: onSecondaryTap,
      onLongPress: onLongPress,
      leading: leading,
      trailing: trailing,
      highlighted: highlighted,
      icon: icon,
      help: help,
      overlay: overlay,
      constraints: constraints,
    );
  }

  const KnownMediaCard(
    this.current, {
    super.key,
    this.onTap,
    this.onDoubleTap,
    this.onSecondaryTap,
    this.onLongPress,
    this.hovered,
    this.leading = const [],
    this.trailing = const [],
    this.highlighted = false,
    this.icon = Icons.play_circle_filled,
    this.help = ds.HelpScope.None,
    this.overlay,
    this.constraints,
  });

  static Widget description(String description) {
    return Flexible(
      child: Tooltip(
        message: description,
        child: Center(
          child: Text(
            description,
            overflow: TextOverflow.ellipsis,
            maxLines: 1,
          ),
        ),
      ),
    );
  }

  static Widget released(String released) {
    return ds.Timestamp.iso8601(
      released,
      format: ds.Timestamp.year,
      inf: ds.Empty,
      neginf: ds.Empty,
    );
  }

  Widget _defaultOverlay(ds.Defaults defaults, api.Known c) {
    return Column(
      mainAxisSize: MainAxisSize.max,
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Expanded(
          child: Padding(
            child: Text(
              c.summary,
              overflow: TextOverflow.ellipsis,
              textAlign: TextAlign.start,
              maxLines: 10,
            ),
            padding: defaults.padding / 2,
          ),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final poster = ds.Image.precache(context, current.image, headers: httpx.localheaders(current.image)) ?? ds.Empty;

    return ConstrainedBox(
      constraints: constraints ?? const BoxConstraints(),
      child: AspectRatio(
        aspectRatio: 2 / 3,
        child: ds.Card(
          SizedBox.expand(
            child: ClipRRect(
              borderRadius: defaults.borderRadius,
              child: ds.Hover(
                poster,
                notifier: hovered,
                overlay: ds.Hover.overlays.icon(
                  context,
                  icon: icon,
                  content: overlay ?? _defaultOverlay(defaults, current),
                ),
              ),
            ),
          ),
          tint: highlighted ? defaults.highlightTint : [],
          alignment: Alignment.center,
          onTap: onTap,
          onDoubleTap: onDoubleTap,
          onSecondaryTap: onSecondaryTap,
          onLongPress: onLongPress,
          help: help,
          leading: [
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: leading.isEmpty ? [KnownMediaCard.description(current.description)] : leading,
            ),
          ],
          trailing: [
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: trailing.isEmpty ? [KnownMediaCard.released(current.released)] : trailing,
            ),
          ],
          fit: FlexFit.tight,
        ),
      ),
    );
  }
}
