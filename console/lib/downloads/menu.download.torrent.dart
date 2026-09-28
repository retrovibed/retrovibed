import 'dart:async';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/design.kit/file.drop.well.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/mimex.dart' as mimex;

Future<List<media.Download>> uploadTorrent(
  BuildContext context, {
  media.FnUploadRequest upload = media.discovered.upload,
}) {
  final tracked = meta.UploadNode.of(context).progress;
  final work = FileDropWell.files(mimetypes: [mimex.bittorrent]).then((evt) {
    return Future.wait(
      evt.files.map((c) {
        final abort = Completer<void>();
        return media.media
            .uploadable(c.path, c.name, c.mimeType!, progress: tracked, abort: abort)
            .then((v) {
              return upload(
                (method, url) => http.AbortableMultipartRequest(method, url, abortTrigger: abort.future)..files.add(v),
              );
            })
            .then<media.DownloadBeginResponse?>((uploaded) {
              return media.discovered.download(
                uploaded.media.id,
                options: [authn.request(authn.AuthzCache.meta(context))],
              );
            })
            .catchError((_) => null, test: httpx.ErrorsTest.aborted);
      }),
    );
  });

  // cancelled torrents resolve to null and drop out rather than queueing an empty download.
  return work.then((v) => v.nonNulls.map((v) => v.download).toList());
}

PopupMenuEntry<String> MenuItemDownloadTorrent(
  BuildContext context,
  Function(Future<List<media.Download>>) onDownload,
) {
  return PopupMenuItem<String>(
    child: ds.LoadingListTile(
      leading: const Icon(Icons.file_download_outlined),
      title: const Text("Download Torrent"),
      onPressed: () => onDownload(
        uploadTorrent(
          context,
        ),
      ).catchError((cause) => debugPrint('upload torrent failed: $cause')),
    ),
  );
}
