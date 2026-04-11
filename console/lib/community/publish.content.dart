import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/library.dart' as lib;

class PublishContent extends StatelessWidget {
  final void Function(media.Download) onSelect;
  final media.FnMediaSearch search;
  final media.FnUploadRequest upload;

  const PublishContent({
    super.key,
    required this.onSelect,
    this.search = media.media.search,
    this.upload = media.media.upload,
  });

  @override
  Widget build(BuildContext context) {
    return lib.AvailableListDisplay(
      search: search,
      upload: upload,
      row:
          (v) => LibraryRow(
            item: media.Download(media: v),
            onTap: () => onSelect(media.Download(media: v)),
          ),
    );
  }
}

class LibraryRow extends StatelessWidget {
  final media.Download item;
  final VoidCallback onTap;

  const LibraryRow({super.key, required this.item, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return ds.TableRow(
      [
        Icon(Icons.folder, size: 32),
        Expanded(
          child: Text(
            item.media.description,
            style: theme.textTheme.bodyLarge,
            overflow: TextOverflow.ellipsis,
          ),
        ),
        Icon(Icons.chevron_right),
      ],
      onTap: onTap,
    );
  }
}
