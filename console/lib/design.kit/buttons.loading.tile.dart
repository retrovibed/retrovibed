import 'package:flutter/material.dart';
import 'buttons.loading.icon.dart' show AsyncVoidCallback;
import 'noop.dart';
import 'stateful.dart';
import 'help.dart';

// A full-width ListTile that shows a spinner in place of [leading] while
// [onPressed] is pending and disables re-tap until it completes.
class LoadingListTile extends StatefulWidget {
  final AsyncVoidCallback onPressed;
  final Widget leading;
  final Widget title;
  final Widget help;
  final bool disabled;

  const LoadingListTile({
    super.key,
    required this.leading,
    required this.title,
    this.help = HelpScope.None,
    this.onPressed = fnAsyncNoop,
    this.disabled = false,
  });

  @override
  State<LoadingListTile> createState() => _LoadingListTileState();
}

class _LoadingListTileState extends State<LoadingListTile> with LoadingState {
  @override
  void initState() {
    super.initState();
    loading = false;
  }

  void _handlePress() {
    if (loading) return;
    setState(() {
      loading = true;
    });
    widget.onPressed().whenComplete(() {
      setState(() {
        loading = false;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    final disabled = loading || widget.disabled;
    return Help(
      ListTile(
        leading: loading
            ? const SizedBox(
                width: 24,
                height: 24,
                child: CircularProgressIndicator(strokeWidth: 2.0),
              )
            : widget.leading,
        title: widget.title,
        enabled: !disabled,
        hoverColor: Colors.transparent,
        onTap: disabled || widget.onPressed == fnAsyncNoop ? null : _handlePress,
      ),
      widget.help,
    );
  }
}
