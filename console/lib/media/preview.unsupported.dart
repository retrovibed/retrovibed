import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/mimex.dart' as mimex;
import './button.play.dart';
import './media.pb.dart';

// what the reader gets when the content cannot be rendered in place: enough to identify
// the file, and the action that fetches it.
class PreviewUnsupported extends StatelessWidget {
  final Media current;
  final Widget description;

  const PreviewUnsupported({
    super.key,
    required this.current,
    this.description = ds.Empty,
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final theme = Theme.of(context);

    return Padding(
      padding: defaults.padding,
      child: Row(
        spacing: defaults.spacing,
        children: [
          Icon(mimex.icon(current.mimetype)),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(current.description, overflow: TextOverflow.ellipsis),
                DefaultTextStyle(
                  style: theme.textTheme.bodySmall!.copyWith(color: theme.hintColor),
                  child: description == ds.Empty ? Text(current.mimetype) : description,
                ),
              ],
            ),
          ),
          ds.LoadingIconButton(
            onPressed: DownloadAction(context, current),
            icon: const Icon(Icons.download),
          ),
        ],
      ),
    );
  }
}
