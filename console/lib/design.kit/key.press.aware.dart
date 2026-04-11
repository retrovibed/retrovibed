import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import './screens.dart' as screens;

class KeyPressAware extends StatefulWidget {
  final Future<void> Function() onPress;
  final Set<LogicalKeyboardKey> keys;
  final Widget child;

  const KeyPressAware(
    this.child, {
    Key? key,
    required this.onPress,
    required this.keys,
  }) : super(key: key);

  /// Factory constructor that presets the keys for common deletion actions.
  factory KeyPressAware.delete(
    Widget child, {
    Key? key,
    required Future<void> Function() onPress,
  }) {
    return KeyPressAware(
      child,
      key: key,
      onPress: onPress,
      keys: {LogicalKeyboardKey.delete, LogicalKeyboardKey.backspace},
    );
  }

  /// Factory constructor that presets the key for refresh (F5).
  factory KeyPressAware.refresh(
    Widget child, {
    Key? key,
    required Future<void> Function() onPress,
  }) {
    return KeyPressAware(
      child,
      key: key,
      onPress: onPress,
      keys: {LogicalKeyboardKey.f5},
    );
  }

  @override
  State<KeyPressAware> createState() => _KeyPressAwareState();
}

class _KeyPressAwareState extends State<KeyPressAware> {
  bool _processing = false;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  KeyEventResult _handleKeyPress(FocusNode node, KeyEvent event) {
    if (event is! KeyDownEvent || _processing) {
      return KeyEventResult.ignored;
    }

    if (!widget.keys.contains(event.logicalKey)) {
      return KeyEventResult.ignored;
    }

    setState(() {
      _processing = true;
    });

    widget.onPress().whenComplete(() {
      setState(() {
        _processing = false;
      });
    });

    return KeyEventResult.handled;
  }

  @override
  Widget build(BuildContext context) {
    return Focus(
      canRequestFocus: false,
      onKeyEvent: _handleKeyPress,
      child: screens.Loading(loading: _processing, widget.child),
    );
  }
}
