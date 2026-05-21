import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/httpx.dart' as httpx;
import './known.media.dropdown.dart';
import './metadata.edit.dart';
import './api.dart' as api;

class MediaSettings extends StatefulWidget {
  final media.Media current;
  final void Function(Future<media.Media> pending, {bool forced, bool autoclose}) onChange;
  final api.FnKnownSearch knownSearch;
  final EdgeInsets? margin;
  final Future<media.MetadataSyncResponse> Function(
    String torrentId,
    media.Media media, {
    List<httpx.Option> options,
  })
  discoveredMetadataSync;
  final Future<media.MediaUpdateResponse> Function(
    String id,
    media.Media media, {
    List<httpx.Option> options,
  })
  libraryMetadataSync;
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

  const MediaSettings({
    super.key,
    required this.current,
    required this.onChange,
    this.margin,
    this.knownSearch = api.known.search,
    this.discoveredMetadataSync = media.discovered.metadatasync,
    this.libraryMetadataSync = media.media.metadatasync,
    this.discoveredUpdate = media.discovered.update,
    this.discoveredReset = media.discovered.reset,
    this.discoveredGet = media.discovered.get,
  });

  @override
  State<MediaSettings> createState() => _MediaSettingsState(current);
}

class _MediaSettingsState extends State<MediaSettings> {
  bool _dirty = false;
  media.Media _modified;
  List<httpx.Option> _authOptions = [];

  _MediaSettingsState(this._modified);

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    // Cache auth context while it's still valid
    _authOptions = [authn.request(authn.AuthzCache.meta(context))];
  }

  @override
  void deactivate() {
    if (_dirty) {
      media.media
          .update(_modified.id, _modified, options: [authn.request(authn.AuthzCache.meta(context))])
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
              onChange: (Future<media.Media> p) {
                p.then((v) {
                  setState(() {
                    _dirty = true;
                    _modified = v;
                  });
                });
              },
            ),
            KnownMediaDropdown(
              current: _modified.knownMediaId,
              search: widget.knownSearch,
              onChange: (known) {
                if (uuidx.isMin(uuidx.fromString(_modified.torrentId))) {
                  return widget.onChange(
                    widget
                        .libraryMetadataSync(
                          _modified.id,
                          _modified..knownMediaId = known?.id ?? uuidx.min(),
                          options: _authOptions,
                        )
                        .then((v) => v.media),
                    forced: true,
                  );
                }

                widget.onChange(
                  widget
                      .discoveredMetadataSync(
                        _modified.torrentId,
                        _modified..knownMediaId = known?.id ?? uuidx.min(),
                        options: _authOptions,
                      )
                      .then((v) => v.media),
                  forced: true,
                );
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
                                onConfirm: () {
                                  widget
                                      .discoveredUpdate(
                                        _modified.torrentId,
                                        download..verifyAt = DateTime.now().toUtc().toIso8601String(),
                                        options: _authOptions,
                                      )
                                      .then((_) => completion.complete())
                                      .catchError((cause) {
                                        completion.completeError(cause);
                                      });
                                },
                                onCancel: completion.complete,
                              ),
                            ),
                        onTap:
                            () => ds.modals.asyncfn(
                              context,
                              (completion) => ds.Confirmation.yesNo(
                                content: Text(
                                  "Are you sure you want to reset ${_modified.description}?",
                                ),
                                onConfirm: () {
                                  widget
                                      .discoveredReset(_modified.torrentId, options: _authOptions)
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
                                onCancel: completion.complete,
                              ),
                            ),
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
