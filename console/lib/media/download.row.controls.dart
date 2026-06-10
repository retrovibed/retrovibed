import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/timex.dart' as timex;
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
    final defaults = ds.Defaults.of(context);
    final pause = () {
      setState(() => disabled = true);
      api.discovered
          .pause(
            widget.current.media.id,
            options: [authn.request(authn.AuthzCache.meta(context))],
          )
          .then((v) {
            widget.onChange?.call(v.download);
            setState(() => disabled = false);
          })
          .catchError((cause) {
            setState(() {
              disabled = false;
              ds.ErrorBoundary.of(context)?.onError(ds.Error.unknown(cause));
            });
          });
    };
    final completed = () {
      widget.onChange?.call(widget.current);
    };
    final tune = () {
      setState(() => disabled = true);
      api.discovered
          .tune(
            widget.current.media.id,
            api.DownloadTuneRequest(peers: []),
            options: [authn.request(authn.AuthzCache.meta(context))],
          )
          .then((v) {
            setState(() => disabled = false);
          })
          .catchError((cause) {
            setState(() {
              disabled = false;
              ds.ErrorBoundary.of(context)?.onError(ds.Error.unknown(cause));
            });
          });
    };

    final cursor = disabled ? SystemMouseCursors.forbidden : SystemMouseCursors.click;
    final isCompleted = timex.iso8601(widget.current.completedAt).isBefore(timex.inf);

    final primaryicon = isCompleted ? Icon(Icons.check, color: defaults.success) : Icon(Icons.pause_circle_outline);
    final primarytap = isCompleted ? completed : pause;

    return Row(
      children: [
        ds.Help(
          IconButton(
            icon: primaryicon,
            mouseCursor: cursor,
            onPressed: disabled ? null : primarytap,
          ),
          ds.Hint(Text(isCompleted ? "mark as processed and remove from active downloads" : "pause the download")),
        ),
        if (authn.developer(context).debug)
          ds.Help(
            IconButton(icon: Icon(Icons.tune), onPressed: disabled ? null : tune),
            ds.Hint(Text("tune peer connections for this download")),
          ),
      ],
    );
  }
}
