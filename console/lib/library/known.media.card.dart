import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
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
  final Widget? trailing;
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
    Widget? trailing,
    bool highlighted = false,
    IconData? icon = Icons.play_circle_filled,
    Widget help = ds.HelpScope.None,
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
          trailing: trailing,
          highlighted: highlighted,
          icon: icon,
          help: help,
          constraints: constraints,
        );
      },
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
    this.trailing,
    this.highlighted = false,
    this.icon = Icons.play_circle_filled,
    this.help = ds.HelpScope.None,
    this.overlay,
    this.constraints,
  });

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
        trailing ?? const SizedBox(),
      ],
    );
  }

  Map<String, String>? _imageheaders(String original) {
    if (original.isEmpty) return null;
    if (!original.startsWith("https://${httpx.host()}")) return null;
    // when we hit a url that matches the current host library we need to add authentication to the request.
    return <String, String>{"Authorization": httpx.auto_bearer_host()};
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final poster = current.image == ""
        ? ds.Empty
        : Image.network(
            current.image,
            headers: _imageheaders(current.image),
            errorBuilder: (context, error, stackTrace) => const Icon(Icons.image_outlined),
          );

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
            Center(
              child: Text(
                current.description,
                overflow: TextOverflow.ellipsis,
                maxLines: 1,
              ),
            ),
          ],
          trailing: [
            trailing ??
                ds.Timestamp.iso8601(
                  current.released,
                  format: ds.Timestamp.year,
                  inf: ds.Empty,
                  neginf: ds.Empty,
                ),
          ],
          fit: FlexFit.tight,
        ),
      ),
    );
  }
}
