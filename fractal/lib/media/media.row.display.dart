import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as ds;
import './media.pb.dart';

class RowDisplay extends StatelessWidget {
  final Media media;
  final Function()? onTap;
  const RowDisplay({super.key, required this.media, this.onTap = null});

  @override
  Widget build(BuildContext context) {
    const gap = SizedBox(width: 10.0);
    final theming = ds.Defaults.of(context);

    return Container(
      padding: theming.padding,
      child: InkWell(
        onTap: onTap,
        child: Row(
          children: [
            const Icon(Icons.movie),
            gap,
            Expanded(
              child: Text(media.description, overflow: TextOverflow.ellipsis),
            ),
            gap,
            Icon(Icons.download),
          ],
        ),
      ),
    );
  }
}
