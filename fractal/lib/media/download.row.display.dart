import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:fractal/media/api.dart' as api;

class DownloadRowDisplay extends StatelessWidget {
  final api.Download current;
  const DownloadRowDisplay({super.key, required this.current});

  @override
  Widget build(BuildContext context) {
    const gap = SizedBox(width: 10.0);

    return Row(
      children: [
        const Icon(Icons.download),
        gap,
        Expanded(
          child: Text(
            current.media.description,
            overflow: TextOverflow.ellipsis,
          ),
        ),
        gap,
        Expanded(
          child: LinearProgressIndicator(
            value:
                current.downloaded.toInt() / math.max(current.bytes.toInt(), 1),
            semanticsLabel: 'Linear progress indicator',
          ),
        ),
        gap,
        IconButton(
          icon: Icon(Icons.delete),
          onPressed: () => print("delete not yet implemented"),
        ),
      ],
    );
  }
}
