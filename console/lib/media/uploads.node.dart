import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/httpx.dart' as httpx;

// used to track ongoing uploads
class UploadNode extends StatefulWidget {
  final Widget child;
  final Duration delay;
  const UploadNode(this.child, {Key? key, this.delay = Durations.extralong4}) : super(key: key);

  static UploadState of(BuildContext context) {
    final scope = context.dependOnInheritedWidgetOfExactType<_UploadScope>();
    assert(scope != null, 'UploadNode.of() called with no UploadNode ancestor');
    return scope!.state;
  }

  @override
  State<UploadNode> createState() => _UploadState();
}

// handed out by UploadNode.of(context): pass `progress` into httpx.uploadable (or a
// wrapper) to report an upload, and read `uploading` to see every upload in flight,
// keyed by id, as of the last progress event. `remove` drops an entry from `uploading`
// (e.g. to dismiss it from the UI); it does not cancel the underlying HTTP request.
class UploadState {
  final StreamSink<httpx.UploadProgress> progress;
  final Map<String, httpx.UploadProgress> uploading;
  final void Function(String id) remove;
  const UploadState(this.progress, this.uploading, this.remove);
}

class _UploadState extends State<UploadNode> {
  final StreamController<httpx.UploadProgress> _progress = StreamController<httpx.UploadProgress>.broadcast();
  final Map<String, httpx.UploadProgress> _uploading = {};

  @override
  void initState() {
    super.initState();
    _progress.stream.listen((event) {
      setState(() {
        _uploading[event.$1] = event;
      });
      if (event.$4 == event.$5) {
        Future.delayed(widget.delay).then((_) {
          setState(() {
            _uploading.removeWhere((id, evt) => evt.$4 == evt.$5);
          });
        });
      }
    });
  }

  @override
  void dispose() {
    _progress.close();
    super.dispose();
  }

  void _remove(String id) {
    setState(() {
      _uploading.remove(id);
    });
  }

  @override
  Widget build(BuildContext context) {
    // a fresh map each build, so updateShouldNotify sees a new reference and
    // actually notifies dependents; _uploading itself is mutated in place.
    return _UploadScope(
      state: UploadState(_progress.sink, Map.unmodifiable(_uploading), _remove),
      child: widget.child,
    );
  }
}

class _UploadScope extends InheritedWidget {
  final UploadState state;
  const _UploadScope({required this.state, required super.child});

  @override
  bool updateShouldNotify(_UploadScope oldWidget) => state.uploading != oldWidget.state.uploading;
}
