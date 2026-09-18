import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'uploads.node.dart';

// an inline strip of the uploads currently in flight, meant to sit alongside
// things like FileDropWell in a toolbar/leading area rather than as its own screen.
class UploadsRow extends StatelessWidget {
  final EdgeInsets? margin;
  final EdgeInsets? padding;
  const UploadsRow({super.key, this.margin, this.padding});

  @override
  Widget build(BuildContext context) {
    final uploads = UploadNode.of(context).uploading.values.toList();
    if (uploads.isEmpty) return ds.Debug.pink(ds.Empty);
    final defaults = ds.Defaults.of(context);

    return ds.Container(
      margin: margin,
      padding: padding,
      Row(
        spacing: defaults.spacing,
        children: uploads.map((u) => Expanded(child: _UploadChip(u))).toList(),
      ),
    );
  }
}

class _UploadChip extends StatelessWidget {
  final httpx.UploadProgress upload;
  const _UploadChip(this.upload);

  // CircularProgressIndicator throws on NaN/Infinity; treat those as indeterminate.
  static double? _fraction(int uploaded, int total) {
    if (total <= 0 || total == httpx.unknownFileSize) return null;
    return (uploaded / total).clamp(0.0, 1.0);
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final (id, name, mimetype, uploaded, total) = upload;
    final fraction = _fraction(uploaded, total);
    final completed = total != httpx.unknownFileSize && uploaded >= total;
    final detail = fraction == null
        ? ds.bytesx(uploaded).toIEC600272Format()
        : "${ds.bytesx(uploaded).toIEC600272().toStringAsFixed(0)} / ${ds.bytesx(total).toIEC600272Format()}";

    return ds.Container(
      key: ValueKey(id),
      padding: const EdgeInsets.symmetric(horizontal: 8.0, vertical: 4.0),
      decoration: BoxDecoration(
        borderRadius: defaults.borderRadius,
        border: Border.all(color: Theme.of(context).dividerColor),
      ),
      Row(
        spacing: defaults.spacing,
        mainAxisSize: MainAxisSize.max,
        children: [
          SizedBox(width: 14.0, height: 14.0, child: CircularProgressIndicator(strokeWidth: 1.0, value: fraction)),
          Icon(mimex.icon(mimetype), size: 16.0),
          Expanded(child: Text(name, overflow: TextOverflow.ellipsis, maxLines: 1)),
          Text(detail, maxLines: 1),
          if (completed)
            ds.Help(
              IconButton(
                icon: const Icon(Icons.close, size: 16.0),
                visualDensity: VisualDensity.compact,
                onPressed: () => UploadNode.of(context).remove(id),
              ),
              ds.Hint(Text("clear the upload from the list")),
            ),
        ],
      ),
    );
  }
}
