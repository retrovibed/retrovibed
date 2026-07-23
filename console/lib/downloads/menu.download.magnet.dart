import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media.dart' as media;
import 'magnet.links.dart';

Future<List<media.Download>> addMagnetLinks(BuildContext context) {
  final completer = Completer<List<media.Download>>();
  ds.modals.push(
    context,
    MagnetDownloads(
      onSubmitted: (magnets) {
        final pending = magnets.map(
          (v) => media.discovered
              .magnet(
                media.MagnetCreateRequest(uri: v),
                options: [authn.request(authn.AuthzCache.meta(context))],
              )
              .then((created) {
                return media.discovered.download(
                  created.download.media.id,
                  options: [authn.request(authn.AuthzCache.meta(context))],
                );
              }),
        );
        return Future.wait(pending, eagerError: true)
            .then((v) {
              ds.modals.of(context)?.reset();
              completer.complete(v.map((v) => v.download).toList());
            })
            .catchError((cause) {
              completer.completeError(cause);
            });
      },
    ),
  );
  return completer.future;
}

PopupMenuEntry<String> MenuItemDownloadMagnet(BuildContext context, Function(Future<List<media.Download>>) onDownload) {
  return PopupMenuItem<String>(
    child: ds.LoadingListTile(
      leading: const Icon(Icons.link),
      title: const Text("Download Magnet"),
      onPressed: () => onDownload(addMagnetLinks(context)),
    ),
  );
}
