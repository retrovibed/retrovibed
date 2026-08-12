import 'package:flutter/material.dart';
import './buttons.loading.icon.dart' show AsyncVoidCallback;
import './stateful.dart';

// A full-width ListTile that shows a spinner in place of [leading] while
// [onPressed] is pending and disables re-tap until it completes.
class LoadingListTile extends StatefulWidget {
  final AsyncVoidCallback onPressed;
  final Widget leading;
  final Widget title;
  final bool disabled;

  const LoadingListTile({
    super.key,
    required this.onPressed,
    required this.leading,
    required this.title,
    this.disabled = false,
  });

  @override
  State<LoadingListTile> createState() => _LoadingListTileState();
}

class _LoadingListTileState extends State<LoadingListTile> with LoadingState {
  bool _isLoading = false;

  void _handlePress() {
    if (_isLoading) return;
    setState(() {
      _isLoading = true;
    });
    widget.onPressed().whenComplete(() {
      setState(() {
        _isLoading = false;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    final disabled = _isLoading || widget.disabled;
    return ListTile(
      leading: _isLoading
          ? const SizedBox(
              width: 24,
              height: 24,
              child: CircularProgressIndicator(strokeWidth: 2.0),
            )
          : widget.leading,
      title: widget.title,
      enabled: !disabled,
      hoverColor: Colors.transparent,
      onTap: disabled ? null : _handlePress,
    );
  }
}
