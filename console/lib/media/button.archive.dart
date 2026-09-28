import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/uuidx.dart' as uuidx;
import './media.pb.dart';
import './button.play.dart';

class ButtonArchive extends StatelessWidget {
  final Media current;
  final void Function(Media upd)? onChange;
  const ButtonArchive({super.key, required this.current, this.onChange});

  @override
  Widget build(BuildContext context) {
    final disabled = authn.AuthzCache.of(context).meta.current.token.archiveUpload.toInt() <= 0;
    final archivable = ds.LoadingIconButton(
      disabled: disabled,
      tooltip: "mark this file to be archived to cloud storage",
      onPressed: ArchiveAction(
        context,
        current,
        then: (v) {
          onChange?.call(v);
          return v;
        },
      ),
      icon: Icon(Icons.upload),
      help: ds.Hint(const Text("upload this file to your cloud archive")),
    );
    final archiving = ds.LoadingIconButton(
      disabled: disabled,
      tooltip: "this file is marked for archival and is awaiting upload, click to cancel",
      onPressed: ArchiveCancelAction(
        context,
        current,
        then: (v) {
          onChange?.call(v);
          return v;
        },
      ),
      icon: Icon(Icons.pending_outlined),
      help: ds.Hint(const Text("file is queued for archival, tap to cancel the upload")),
    );
    final purge = ds.LoadingIconButton(
      disabled: disabled,
      tooltip: "purge content from your archive",
      onPressed: ArchivePurgeAction(
        context,
        current,
        then: (v) {
          final upd = current..archiveId = uuidx.min();
          onChange?.call(upd);
          return upd;
        },
      ),
      icon: Icon(Icons.delete_forever),
      help: ds.Hint(const Text("remove this file's copy from your cloud archive")),
    );

    return uuidx.pattern(current.archiveId, archivable, archiving, purge);
  }
}
