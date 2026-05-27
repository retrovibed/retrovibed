import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/httpx.dart' as httpx;
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
        background: theme.colorScheme.surfaceContainerLow,
        Column(
          mainAxisAlignment: MainAxisAlignment.start,
          mainAxisSize: MainAxisSize.max,
          spacing: defaults.spacing,
          children: [
            MediaEdit(
              current: _modified,
              padding: defaults.padding,
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
            if (!uuidx.isMinMax(uuidx.fromString(_modified.torrentId)))
              ds.Container(
                background: theme.colorScheme.surface,
                ds.Accordion(
                  expanded: true,
                  description: Text("source details"),
                  content: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      media.DownloadDisplay.fromID(
                        _modified.torrentId,
                        get: widget.discoveredGet,
                        onVerify:
                            (download) => ds.modals.asyncfn(
                              context,
                              (completion) => ds.Confirmation.yesNo(
                                content: Text(
                                  "Are you sure you want to verify ${_modified.description}?",
                                ),
                                onConfirm: (ctx) {
                                  widget
                                      .discoveredUpdate(
                                        _modified.torrentId,
                                        download..verifyAt = DateTime.now().toUtc().toIso8601String(),
                                        options: [authn.request(authn.AuthzCache.meta(ctx))],
                                      )
                                      .then((_) => completion.complete())
                                      .catchError((cause) {
                                        completion.completeError(cause);
                                      });
                                },
                                onCancel: (_) => completion.complete(),
                              ),
                            ),
                        onTap:
                            () => ds.modals.asyncfn(context, (completion) {
                              return ds.Confirmation.yesNo(
                                content: Text(
                                  "Are you sure you want to reset ${_modified.description}?",
                                ),
                                onConfirm: (ctx) {
                                  widget
                                      .discoveredReset(
                                        _modified.torrentId,
                                        options: [authn.request(authn.AuthzCache.meta(ctx))],
                                      )
                                      .then((v) {
                                        widget.onChange(
                                          Future.value(_modified),
                                          forced: true,
                                        );
                                        completion.complete();
                                      })
                                      .catchError((cause) {
                                        completion.completeError(cause);
                                      });
                                },
                                onCancel: (_) => completion.complete(),
                              );
                            }),
                      ),
                    ],
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }
}
