import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/file.drop.well.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;

Future<void> uploadFiles(
  BuildContext context,
  ValueNotifier<media.MediaSearchState> search, {
  media.FnUploadRequest apiupload = media.media.upload,
  List<String> mimetypes = const [],
}) {
  return FileDropWell.pickFiles(mimetypes: mimetypes).then((evt) {
    return Future.wait(
      evt.files.map((c) {
        return media.media.uploadable(c.path, c.name, c.mimeType!).then((v) {
          return apiupload((req) {
            req..files.add(v);
            return req;
          });
        });
      }),
    ).then((_) {
      final freshNext = search.value.next.clone();
      search.value = media.MediaSearchState(next: freshNext, count: search.value.count);
    });
  });
}

PopupMenuEntry<String> MenuItemUploadFiles(
  BuildContext context,
  ValueNotifier<media.MediaSearchState> search, {
  media.FnUploadRequest apiupload = media.media.upload,
  List<String> mimetypes = const [],
}) {
  return PopupMenuItem<String>(
    child: ds.LoadingListTile(
      leading: const Icon(Icons.file_upload_outlined),
      title: const Text("Upload files"),
      onPressed: () => uploadFiles(
        context,
        search,
        apiupload: apiupload,
        mimetypes: mimetypes,
      ).catchError((cause) => debugPrint('upload files failed: $cause')),
    ),
  );
}
