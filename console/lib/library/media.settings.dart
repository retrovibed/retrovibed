import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/httpx.dart' as httpx;
import 'media.download.accordian.dart';
import 'known.media.accordian.dart';
import 'metadata.edit.dart';
import 'api.dart' as api;

class MediaSettings extends StatefulWidget {
  final media.Media current;
  final void Function(Future<media.Media> pending, {bool forced, bool autoclose}) onChange;
  final api.FnKnownSearch knownSearch;
  final EdgeInsets? margin;
  final Future<media.DownloadUpdateResponse> Function(
    String id,
    media.Download download, {
    List<httpx.Option> options,
  })
  discoveredUpdate;
  final Future<media.DownloadDeleteResponse> Function(
    String id, {
    List<httpx.Option> options,
  })
  discoveredReset;
  final Future<media.DownloadMetadataResponse> Function(
    String id, {
    List<httpx.Option> options,
  })
  discoveredGet;
  final Future<media.MediaUpdateResponse> Function(
    String id,
    media.Media upd, {
    List<httpx.Option> options,
  })
  mediaUpdate;
  final Future<media.MediaDeleteResponse> Function(
    String id, {
    List<httpx.Option> options,
  })
  mediaDelete;

  const MediaSettings({
    super.key,
    required this.current,
    required this.onChange,
    this.margin,
    this.knownSearch = api.known.search,
    this.discoveredUpdate = media.discovered.update,
    this.discoveredReset = media.discovered.reset,
    this.discoveredGet = media.discovered.get,
    this.mediaUpdate = media.media.update,
    this.mediaDelete = media.media.delete,
  });

  @override
  State<MediaSettings> createState() => _MediaSettingsState(current);
}

class _MediaSettingsState extends State<MediaSettings> {
  bool _dirty = false;
  media.Media _modified;

  _MediaSettingsState(this._modified);

  @override
  void deactivate() {
    if (_dirty) {
      widget
          .mediaUpdate(_modified.id, _modified, options: [authn.request(authn.AuthzCache.meta(context))])
          .then((v) => widget.onChange(Future.value(v.media)));
    }

    super.deactivate();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);
    return SelectionArea(
      child: ds.Container(
        padding: defaults.padding,
        margin: widget.margin ?? defaults.margin,
        decoration: BoxDecoration(color: theme.colorScheme.surfaceContainerLow),
        Column(
          mainAxisAlignment: MainAxisAlignment.start,
          mainAxisSize: MainAxisSize.max,
          spacing: defaults.spacing,
          children: [
            MediaEdit(
              current: _modified,
              padding: defaults.padding,
              deletable: ds.LoadingIconButton(
                icon: Icon(Icons.delete_forever_rounded),
                onPressed: () => ds.modals.asyncfn(
                  context,
                  ds.Confirmation.dangerous(
                    content: Text(
                      "Are you sure you want to delete ${_modified.description}?",
                    ),
                    onConfirm: (ctx) => httpx
                        .withRetry(
                          () => widget.mediaDelete(
                            _modified.id,
                            options: [authn.request(authn.AuthzCache.meta(ctx))],
                          ),
                        )
                        .then((_) {
                          widget.onChange(Future.value(_modified), forced: true, autoclose: true);
                        }),
                  ),
                ),
              ),
              closable: ds.LoadingIconButton.close(
                onPressed: () async => widget.onChange(Future.value(widget.current), autoclose: true),
              ),
              onChange: (Future<media.Media> p) {
                p.then((v) {
                  setState(() {
                    _dirty = true;
                    _modified = v;
                  });
                });
              },
            ),
            KnownMediaAccordian(
              api.known.autodetect(_modified, options: [authn.request(authn.AuthzCache.meta(context))]),
            ),
            if (!uuidx.isMinMax(uuidx.fromString(_modified.torrentId)))
              MediaDownloadAccordian(
                torrentId: _modified.torrentId,
                description: _modified.description,
                discoveredGet: widget.discoveredGet,
                discoveredUpdate: widget.discoveredUpdate,
                discoveredReset: widget.discoveredReset,
                onReset: () => widget.onChange(Future.value(_modified), forced: true),
              ),
          ],
        ),
      ),
    );
  }
}
