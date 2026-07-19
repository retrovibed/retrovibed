import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media.dart' as media;
import 'magnet.links.dart';

Future<void> addMagnetLinks(BuildContext context) {
  final completer = Completer<void>();
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
            .then((_) {
              ds.modals.of(context)?.reset();
              completer.complete();
            })
            .catchError((cause) {
              completer.completeError(cause);
            });
      },
    ),
  );
  return completer.future;
}

PopupMenuEntry<String> MenuItemDownloadMagnet(BuildContext context) {
  return PopupMenuItem<String>(
    child: ds.LoadingListTile(
      leading: const Icon(Icons.link),
      title: const Text("Download Magnet"),
      onPressed: () => addMagnetLinks(context).catchError((cause) => debugPrint('add magnet link failed: $cause')),
    ),
  );
}
