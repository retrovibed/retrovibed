import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/library/known.media.card.dart';
import 'package:retrovibed/media/media.pb.dart';
import 'package:retrovibed/media.dart' as _media;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'api.dart' as api;
import 'known.media.source.dart';

class KnownMediaDisplay extends StatefulWidget {
  final Future<api.Known> pending;
  final _media.Media media;
  final GestureTapCallback? onTap;
  final GestureTapCallback? onDoubleTap;
  final void Function()? onSettings;
  final void Function(_media.Media upd)? onChange;
  final List<Widget> trailing;
  final bool highlighted;
  final Widget help;

  const KnownMediaDisplay(
    this.pending, {
    super.key,
    this.onTap,
    this.onDoubleTap,
    this.onSettings,
    this.onChange,
    this.trailing = const [],
    this.highlighted = false,
    this.help = ds.HelpScope.None,
    required this.media,
  });

  factory KnownMediaDisplay.missing(
    Media m, {
    Key? key,
    GestureTapCallback? onTap,
    GestureTapCallback? onDoubleTap,
    void Function()? onSettings,
    void Function(_media.Media upd)? onChange,
    List<Widget> trailing = const [],
    bool highlighted = false,
    Widget help = ds.HelpScope.None,
  }) {
    return KnownMediaDisplay(
      api.known.autodetect(m),
      media: m,
      key: key ?? ValueKey(m.id),
      onTap: onTap,
      onDoubleTap: onDoubleTap,
      onSettings: onSettings,
      onChange: onChange,
      trailing: trailing,
      highlighted: highlighted,
      help: help,
    );
  }

  factory KnownMediaDisplay.auto(
    BuildContext context,
    Media m, {
    Key? key,
    GestureTapCallback? onTap,
    GestureTapCallback? onDoubleTap,
    void Function()? onSettings,
    void Function(_media.Media upd)? onChange,
    List<Widget> trailing = const [],
    bool highlighted = false,
    Widget help = ds.HelpScope.None,
  }) {
    final resolvedKey = key ?? ValueKey(uuidx.md5x("${m.id}.${m.updatedAt}"));

    if (uuidx.isMinMax(uuidx.fromString(m.knownMediaId))) {
      return KnownMediaDisplay.missing(
        m,
        key: resolvedKey,
        onTap: onTap,
        onDoubleTap: onDoubleTap,
        onSettings: onSettings,
        onChange: onChange,
        trailing: trailing,
        highlighted: highlighted,
        help: help,
      );
    }

    final authz = authn.AuthzCache.meta(context);
    return KnownMediaDisplay(
      api.known
          .cached(
            m.knownMediaId,
            () => api.known.get(
              m.knownMediaId,
              options: [authn.request(authz)],
            ),
          )
          .then((w) => (w.known..description = m.description)),
      media: m,
      key: resolvedKey,
      onTap: onTap,
      onDoubleTap: onDoubleTap,
      onSettings: onSettings,
      onChange: onChange,
      trailing: trailing,
      highlighted: highlighted,
      help: help,
    );
  }

  static const hintPlayMedia = ds.Hint(
    Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text("tap — play media item and auto-generate a playlist from results"),
        Text("long press — show or hide additional controls"),
        Text("settings — edit metadata, tags, and playback options"),
      ],
    ),
  );

  static const hintReleases = ds.Hint(
    Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text("Newly released media discovered through content partnerships."),
        Text("Hover over the card to reveal a summary and the download icon."),
        Text("Tap the card to add it to your library."),
      ],
    ),
  );

  static const hintRecommendations = ds.Hint(
    Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text("Suggested content for you."),
        Text("Hover over a card to reveal a summary and the download icon."),
        Text("Tap the card to add the item to your library."),
        Text("Long Tap the card to clear the recommendation"),
        Text("Use the refresh button to generate a new set of recommendations."),
      ],
    ),
  );

  static _KnownMediaDisplayState? of(BuildContext context) {
    return context.findAncestorStateOfType<_KnownMediaDisplayState>();
  }

  @override
  State<StatefulWidget> createState() => _KnownMediaDisplayState();
}

class _KnownMediaDisplayState extends State<KnownMediaDisplay> {
  final hovered = ValueNotifier(false);

  api.Known current = api.Known(
    id: "",
    description: "",
    summary: "",
    rating: 0.0,
    image: "",
  );

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  void initState() {
    super.initState();
    widget.pending.then((v) {
      setState(() {
        current = v;
      });
    });
  }

  @override
  void dispose() {
    hovered.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final archivable = ds.LoadingIconButton(
      tooltip: "mark this file to be archived to cloud storage",
      onPressed: _media.ArchiveAction(
        context,
        widget.media,
        then: (v) {
          widget.onChange?.call(v);
          return v;
        },
      ),
      icon: Icon(Icons.upload),
    );
    final archiving = ds.LoadingIconButton(
      tooltip: "this file is marked for archival and is awaiting upload, click to cancel",
      onPressed: _media.ArchiveCancelAction(
        context,
        widget.media,
        then: (v) {
          widget.onChange?.call(v);
          return v;
        },
      ),
      icon: Icon(Icons.pending_outlined),
    );
    final purge = ds.LoadingIconButton(
      tooltip: "purge content from your archive",
      onPressed: _media.ArchivePurgeAction(
        context,
        widget.media,
        then: (v) {
          final upd = widget.media..archiveId = uuidx.min();
          widget.onChange?.call(upd);
          return upd;
        },
      ),
      icon: Icon(Icons.delete_forever),
    );

    return KnownMediaCard(
      current,
      highlighted: widget.highlighted,
      hovered: hovered,
      help: widget.help,
      onTap: widget.onTap,
      onDoubleTap: widget.onDoubleTap,
      onLongPress: () {
        setState(() => hovered.value = !hovered.value);
      },
      icon: mimex.icon(widget.media.mimetype),
      trailing: [
        ds.layout(
          (context, constraints) {
            final defaults = ds.Defaults.of(context);
            return Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                if (constraints.maxWidth >= 260) ds.Rating(rating: current.rating),
                Flexible(child: KnownMediaSource(current)),
                Visibility(
                  visible: (authn.AuthzCache.of(context).meta.current.token.archiveUpload.toInt()) > 0,
                  child: uuidx.pattern(widget.media.archiveId, archivable, archiving, purge),
                ),
                ds.LoadingIconButton(
                  tooltip: "download this file to your downloads folder",
                  onPressed: _media.DownloadAction(context, widget.media),
                  icon: Icon(Icons.download),
                ),
                ds.LoadingIconButton(
                  tooltip: "file information and management",
                  onPressed: ds.LoadingIconButton.convert(widget.onSettings),
                  icon: Icon(Icons.tune),
                ),
                if (defaults.mobile)
                  ds.LoadingIconButton.info(
                    tooltip: "show media details",
                    toggled: hovered.value,
                    onPressed: () async => setState(() => hovered.value = !hovered.value),
                  ),
                ...widget.trailing,
              ],
            );
          },
        ),
      ],
    );
  }
}
