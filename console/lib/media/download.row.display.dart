import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:console/designkit.dart' as ds;
import './api.dart' as api;

class DownloadRowDisplay extends StatelessWidget {
  final api.Download current;
  final Widget? Function(BuildContext)? trailing;
  const DownloadRowDisplay({super.key, required this.current, this.trailing});

  @override
  Widget build(BuildContext context) {
    const gap = SizedBox(width: 10.0);

    return ds.ErrorBoundary(
      Row(
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
          Text(current.peers.toString()),
          gap,
          Expanded(
            child: LinearProgressIndicator(
              value:
                  current.downloaded.toInt() /
                  math.max(current.bytes.toInt(), 1),
              semanticsLabel: 'Linear progress indicator',
            ),
          ),
          gap,
          trailing?.call(context) ?? const SizedBox(),
        ],
      ),
    );
  }
}
