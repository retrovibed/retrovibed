import 'package:flutter/material.dart';
import 'package:console/designkit.dart' as ds;
import './api.dart' as api;

class DownloadRowControls extends StatefulWidget {
  final api.Download current;
  final Function(api.Download)? onChange;
  const DownloadRowControls({super.key, required this.current, this.onChange});

  @override
  State<DownloadRowControls> createState() => _ControlState();
}

class _ControlState extends State<DownloadRowControls> {
  bool disabled = false;

  @override
  Widget build(BuildContext context) {
    final cursor =
        disabled ? SystemMouseCursors.forbidden : SystemMouseCursors.click;

    return Row(
      children: [
        IconButton(
          icon: Icon(Icons.pause_circle_outline),
          mouseCursor: cursor,
          onPressed:
              disabled
                  ? null
                  : () {
                    setState(() => disabled = true);
                    api.discovered
                        .pause(widget.current.media.id)
                        .then((v) {
                          widget.onChange?.call(v.download);
                          setState(() => disabled = false);
                        })
                        .catchError((cause) {
                          setState(() {
                            disabled = false;
                            ds.ErrorBoundary.of(
                              context,
                            )?.onError(ds.Error.unknown(cause));
                          });
                        });
                  },
        ),
      ],
    );
  }
}
