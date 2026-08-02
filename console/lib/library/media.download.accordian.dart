import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/httpx.dart' as httpx;

class MediaDownloadAccordian extends StatelessWidget {
  final String torrentId;
  final String description;
  final Future<media.DownloadUpdateResponse> Function(
    String id,
    media.Download download, {
    List<httpx.Option> options,
  }) discoveredUpdate;
  final Future<media.DownloadDeleteResponse> Function(
    String id, {
    List<httpx.Option> options,
  }) discoveredReset;
  final Future<media.DownloadMetadataResponse> Function(
    String id, {
    List<httpx.Option> options,
  }) discoveredGet;
  final void Function()? onReset;

  const MediaDownloadAccordian({
    super.key,
    required this.torrentId,
    required this.description,
    this.discoveredUpdate = media.discovered.update,
    this.discoveredReset = media.discovered.reset,
    this.discoveredGet = media.discovered.get,
    this.onReset,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return ds.Container(
      decoration: BoxDecoration(color: theme.colorScheme.surface),
      ds.Accordion(
        expanded: true,
        description: Text("source details"),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            media.DownloadDisplay.fromID(
              torrentId,
              get: discoveredGet,
              onVerify:
                  (download) => ds.modals.asyncfn(
                    context,
                    (completion) => ds.Confirmation.yesNo(
                      content: Text(
                        "Are you sure you want to verify $description?",
                      ),
                      onConfirm: (ctx) {
                        discoveredUpdate(
                          torrentId,
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
                        "Are you sure you want to reset $description?",
                      ),
                      onConfirm: (ctx) {
                        discoveredReset(
                          torrentId,
                          options: [authn.request(authn.AuthzCache.meta(ctx))],
                        )
                            .then((_) {
                              onReset?.call();
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
    );
  }
}
