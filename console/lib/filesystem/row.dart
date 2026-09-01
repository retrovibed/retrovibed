import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/library/metadata.icons.dart' as icons;

class FilesystemRow extends StatelessWidget {
  final media.Media current;
  final Future<void> Function()? onTap;
  final List<Widget> trailing;
  final bool highlighted;

  const FilesystemRow({
    super.key,
    required this.current,
    this.onTap,
    this.trailing = const [],
    this.highlighted = false,
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final directory = current.mimetype == mimex.directory;

    return media.RowDisplay(
      media: current,
      highlighted: highlighted,
      leading: [Icon(mimex.icon(current.mimetype))],
      // a directory has no bytes to archive, so the local/pending/archived indicator would
      // report every directory as local only.
      trailing: [
        ...trailing,
        Visibility(
          visible: !directory,
          child: icons.archived(current.archiveId, defaults: defaults),
        ),
      ],
      onTap: onTap,
    );
  }
}
