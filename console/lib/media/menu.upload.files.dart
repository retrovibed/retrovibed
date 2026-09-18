import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/file.drop.well.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/mimex.dart' as mimex;
import 'api.dart' as api;
import 'search.state.dart' show MediaSearchState;

Future<void> uploadfiles(
  BuildContext context,
  ValueNotifier<MediaSearchState> search, {
  api.FnUploadRequest apiupload = api.media.upload,
  List<String> mimetypes = const [],
}) {
  final progress = meta.UploadNode.of(context).progress;
  return FileDropWell.files(mimetypes: mimetypes).then((evt) {
    return Future.wait(
      evt.files.map((c) {
        return api.media.uploadable(c.path, c.name, c.mimeType!, progress: progress).then((v) {
          return apiupload((req) {
            req..files.add(v);
            return req;
          });
        });
      }),
    ).then((_) {
      final freshNext = search.value.next.clone();
      search.value = MediaSearchState(next: freshNext, count: search.value.count);
    });
  });
}

PopupMenuEntry<String> MenuItemUploadFiles(
  BuildContext context,
  ValueNotifier<MediaSearchState> search, {
  api.FnUploadRequest apiupload = api.media.upload,
}) {
  return PopupMenuItem<String>(
    child: ValueListenableBuilder<MediaSearchState>(
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
