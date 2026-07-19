import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/file.drop.well.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;

Future<void> uploadTorrent(
  BuildContext context, {
  media.FnUploadRequest upload = media.discovered.upload,
}) {
  return FileDropWell.pickFiles(mimetypes: [mimex.bittorrent]).then((evt) {
    return Future.wait(
      evt.files.map((c) {
        return media.media.uploadable(c.path, c.name, c.mimeType!).then((v) {
          return upload((req) {
            req..files.add(v);
            return req;
          }).then((uploaded) {
            return media.discovered.download(
              uploaded.media.id,
              options: [authn.request(authn.AuthzCache.meta(context))],
            );
          });
        });
      }),
    );
  });
}

PopupMenuEntry<String> MenuItemDownloadTorrent(BuildContext context) {
  return PopupMenuItem<String>(
    child: ds.LoadingListTile(
      leading: const Icon(Icons.file_download_outlined),
      title: const Text("Download Torrent"),
      onPressed: () => uploadTorrent(context).catchError((cause) => debugPrint('upload torrent failed: $cause')),
    ),
  );
}
