import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/file.drop.well.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;

Future<void> uploadfiles(
  BuildContext context,
  ValueNotifier<media.MediaSearchState> search, {
  media.FnUploadRequest apiupload = media.media.upload,
  List<String> mimetypes = const [],
}) {
  return FileDropWell.files(mimetypes: mimetypes).then((evt) {
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
}) {
  return PopupMenuItem<String>(
    child: ValueListenableBuilder<media.MediaSearchState>(
      valueListenable: search,
      builder: (context, state, _) => ds.LoadingListTile(
        leading: const Icon(Icons.file_upload_outlined),
        title: Text("Upload ${mimex.CategoryOptionsLabel.text(state.next.mimetypes)}"),
        onPressed: () => uploadfiles(
          context,
          search,
          apiupload: apiupload,
          mimetypes: state.next.mimetypes,
        ),
      ),
    ),
  );
}
