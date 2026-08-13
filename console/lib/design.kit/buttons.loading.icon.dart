import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/help.dart';

// The original nullable synchronous function type
// typedef _NullableVoidCallback = void Function()?;

// The target non-nullable asynchronous function type
typedef AsyncVoidCallback = Future<void> Function();

// The custom button that handles the asynchronous state
class LoadingIconButton extends StatefulWidget {
  static const double _defaultIconSize = 24.0;
  final AsyncVoidCallback onPressed;
  final Widget icon;
  final String? tooltip;
  final bool disabled;
  final Widget help;
  final double? value;
  final double iconSize;
  final FocusNode? focusNode;

  final bool? toggled;

  const LoadingIconButton({
    super.key,
    required this.onPressed,
    required this.icon,
    this.tooltip,
    this.value,
    this.iconSize = _defaultIconSize,
    this.disabled = false,
    this.toggled,
    this.help = HelpScope.None,
    this.focusNode,
  });

  factory LoadingIconButton.create({
    required AsyncVoidCallback onPressed,
    Key? key,
    String? tooltip,
    bool disabled = false,
    bool? toggled,
    double iconSize = _defaultIconSize,
    double? value,
    Widget help = HelpScope.None,
  }) {
    return LoadingIconButton(
      key: key,
      onPressed: onPressed,
      icon: const Icon(Icons.add),
      tooltip: tooltip,
      disabled: disabled,
      toggled: toggled,
      iconSize: iconSize,
      value: value,
      help: help,
    );
  }

  factory LoadingIconButton.minimize({
    required AsyncVoidCallback onPressed,
    Key? key,
    String? tooltip,
    bool disabled = false,
    bool? toggled,
    double iconSize = _defaultIconSize,
    double? value,
    Widget help = HelpScope.None,
  }) {
    return LoadingIconButton(
      key: key,
      onPressed: onPressed,
      icon: const Icon(Icons.minimize),
      tooltip: tooltip,
      disabled: disabled,
      toggled: toggled,
      iconSize: iconSize,
      value: value,
      help: help,
    );
  }

  factory LoadingIconButton.logout({
    required AsyncVoidCallback onPressed,
    Key? key,
    String? tooltip,
    bool disabled = false,
    bool? toggled,
    double iconSize = _defaultIconSize,
    double? value,
    Widget help = HelpScope.None,
  }) {
    return LoadingIconButton(
      key: key,
      onPressed: onPressed,
      icon: const Icon(Icons.logout),
      tooltip: tooltip,
      disabled: disabled,
      toggled: toggled,
      iconSize: iconSize,
      value: value,
      help: help,
    );
  }

  factory LoadingIconButton.close({
    required AsyncVoidCallback onPressed,
    Key? key,
    String? tooltip,
    bool disabled = false,
    bool? toggled,
    double iconSize = _defaultIconSize,
    double? value,
    Widget help = HelpScope.None,
  }) {
    return LoadingIconButton(
      key: key,
      onPressed: onPressed,
      icon: const Icon(Icons.close),
      tooltip: tooltip,
      disabled: disabled,
      toggled: toggled,
      iconSize: iconSize,
      value: value,
      help: help,
    );
  }

  factory LoadingIconButton.edit({
    required AsyncVoidCallback onPressed,
    Key? key,
    String? tooltip,
    bool disabled = false,
    bool? toggled,
    double iconSize = _defaultIconSize,
    double? value,
    Widget help = HelpScope.None,
  }) {
    return LoadingIconButton(
      key: key,
      onPressed: onPressed,
      icon: const Icon(Icons.edit),
      tooltip: tooltip,
      disabled: disabled,
      toggled: toggled,
      iconSize: iconSize,
      value: value,
      help: help,
    );
  }

  factory LoadingIconButton.delete({
    required AsyncVoidCallback onPressed,
    Key? key,
    String? tooltip,
    bool disabled = false,
    bool? toggled,
    double iconSize = _defaultIconSize,
    double? value,
    Widget help = HelpScope.None,
  }) {
    return LoadingIconButton(
      key: key,
      onPressed: onPressed,
      icon: const Icon(Icons.delete),
      tooltip: tooltip,
      disabled: disabled,
      toggled: toggled,
      iconSize: iconSize,
      value: value,
      help: help,
    );
  }

  factory LoadingIconButton.refresh({
    required AsyncVoidCallback onPressed,
    Key? key,
    String? tooltip,
    bool disabled = false,
    bool? toggled,
    double iconSize = _defaultIconSize,
    double? value,
    Widget help = HelpScope.None,
  }) {
    return LoadingIconButton(
      key: key,
      onPressed: onPressed,
      icon: const Icon(Icons.refresh),
      tooltip: tooltip,
      disabled: disabled,
      toggled: toggled,
      iconSize: iconSize,
      value: value,
      help: help,
    );
  }

  factory LoadingIconButton.search({
    required AsyncVoidCallback onPressed,
    Key? key,
    String? tooltip,
    bool disabled = false,
    bool? toggled,
    double iconSize = _defaultIconSize,
    double? value,
    Widget help = HelpScope.None,
  }) {
    return LoadingIconButton(
      key: key,
      onPressed: onPressed,
      icon: const Icon(Icons.search),
      tooltip: tooltip,
      disabled: disabled,
      toggled: toggled,
      iconSize: iconSize,
      value: value,
      help: help,
    );
  }

  factory LoadingIconButton.queue({
    required AsyncVoidCallback onPressed,
    Key? key,
    String? tooltip,
    bool disabled = false,
    bool? toggled,
    double iconSize = _defaultIconSize,
    double? value,
    Widget help = HelpScope.None,
  }) {
    return LoadingIconButton(
      key: key,
      onPressed: onPressed,
      icon: const Icon(Icons.queue_music),
      tooltip: tooltip,
      disabled: disabled,
      toggled: toggled,
      iconSize: iconSize,
      value: value,
      help: help,
    );
  }

  factory LoadingIconButton.info({
    required AsyncVoidCallback onPressed,
    Key? key,
    String? tooltip,
    bool disabled = false,
    bool? toggled,
    double iconSize = _defaultIconSize,
    double? value,
    Widget help = HelpScope.None,
  }) {
    return LoadingIconButton(
      key: key,
      onPressed: onPressed,
      icon: const Icon(Icons.info_outline_rounded),
      tooltip: tooltip,
      disabled: disabled,
      toggled: toggled,
      iconSize: iconSize,
      value: value,
      help: help,
    );
  }

  static AsyncVoidCallback convert(VoidCallback? fn) {
    return () async {
      fn?.call();
    };
  }

  @override
  _LoadingIconButtonState createState() => _LoadingIconButtonState();
}

class _LoadingIconButtonState extends State<LoadingIconButton> {
  bool _isLoading = false;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _handlePress() {
    if (_isLoading) {
      return;
    }

    // Update state to show loading spinner and disable the button
    setState(() {
      _isLoading = true;
    });

    // Execute the user's asynchronous function and await its completion
    widget
        .onPressed()
        .catchError((cause) {
          debugPrint('$cause');
          if (!mounted) return;
          // Handle any errors that occur during the asynchronous operation
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text("An error occurred: $cause")),
          );
        })
        .whenComplete(() {
          setState(() {
            _isLoading = false;
          });
        });
  }

  // CircularProgressIndicator throws on NaN/Infinity; treat those as indeterminate.
  static double? _sanitize(double? v) {
    if (v == null || v.isNaN || v.isInfinite) return null;
    return v.clamp(0.0, 1.0);
  }

  @override
  Widget build(BuildContext context) {
    final isOn = widget.toggled == true;
    final color = isOn ? Theme.of(context).colorScheme.primary : null;
    return Help(
      IconButton(
        tooltip: widget.tooltip,
        mouseCursor: _isLoading || widget.disabled ? SystemMouseCursors.basic : SystemMouseCursors.click,
        onPressed: _isLoading || widget.disabled ? null : _handlePress,
        focusNode: widget.focusNode,
        iconSize: widget.iconSize,
        style: IconButton.styleFrom(
          padding: EdgeInsets.all(widget.iconSize * 0.25),
          minimumSize: Size.zero,
          tapTargetSize: MaterialTapTargetSize.shrinkWrap,
        ),
        color: color,
        icon: _isLoading
            ? SizedBox(
                width: widget.iconSize,
                height: widget.iconSize,
                child: CircularProgressIndicator(strokeWidth: 2.0, value: _sanitize(widget.value)),
              )
            : widget.icon,
      ),
      widget.help,
    );
  }
}
