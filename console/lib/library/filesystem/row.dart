import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import '../metadata.icons.dart' as icons;

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
    final folder = current.mimetype == mimex.directory;

    return media.RowDisplay(
      media: current,
      highlighted: highlighted,
      leading: [Icon(mimex.icon(current.mimetype))],
      // a folder has no bytes to archive, so the local/pending/archived indicator would
      // report every folder as local only.
      trailing: [
        ...trailing,
        Visibility(
          visible: !folder,
          child: icons.archived(current.archiveId, defaults: defaults),
        ),
      ],
      onTap: onTap,
    );
  }
}
