import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as ds;
import 'package:desktop_drop/desktop_drop.dart';

class FileDropWell extends StatefulWidget {
  final Widget child;
  final Function()? onTap;
  final Future<Widget?> Function(DropDoneDetails i) onDropped;
  const FileDropWell(
    this.onDropped, {
    super.key,
    this.child = const Center(
      child: Column(
        mainAxisSize: MainAxisSize.max,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.filter_rounded),
          SelectableText("Drop Files to add them to your library."),
        ],
      ),
    ),
    this.onTap,
  });

  @override
  _FileDropWell createState() => _FileDropWell();
}

class _FileDropWell extends State<FileDropWell> {
  bool _dragging = false;
  bool _loading = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return DropTarget(
      onDragDone: (evt) {
        setState(() {
          _loading = true;
        });
        widget.onDropped(evt).then((w) {}).whenComplete(() {
          setState(() {
            _loading = false;
          });
        });
      },
      onDragEntered: (detail) {
        setState(() {
          _dragging = true;
        });
      },
      onDragExited: (detail) {
        setState(() {
          _dragging = false;
        });
      },
      child: Container(
        color: _dragging ? theme.highlightColor : null,
        child: ds.Loading(
          loading: _loading,
          child: Row(
            mainAxisSize: MainAxisSize.max,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [widget.child],
          ),
        ),
      ),
    );
  }
}
