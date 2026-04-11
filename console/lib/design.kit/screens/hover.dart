import 'package:flutter/material.dart' as material;
import '../theme.defaults.dart';
import './overlay.dart';

class _Overlays {
  const _Overlays();

  material.Widget icon(
    material.BuildContext context, {
    required material.Widget content,
    material.IconData? icon = material.Icons.play_circle_filled,
  }) {
    final defaults = Defaults.of(context);
    return material.Container(
      child: material.Stack(
        children: [
          content,
          material.Center(
            child: material.Icon(
              icon,
              size: 64,
              color: (material.Theme.of(context).iconTheme.color ?? material.Colors.white).withValues(
                alpha: defaults.highlight.a,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class Hover extends material.StatefulWidget {
  static const overlays = _Overlays();
  static const _zero = const material.SizedBox();
  final material.Widget child;
  final material.Widget overlay;
  final material.ValueNotifier<bool>? notifier;

  const Hover(this.child, {super.key, required this.overlay, this.notifier});

  @override
  material.State<Hover> createState() => _HoverState();
}

class _HoverState extends material.State<Hover> {
  final _hovered = material.ValueNotifier(false);

  bool get _active => _hovered.value || (widget.notifier?.value ?? false);

  @override
  void initState() {
    super.initState();
    _hovered.addListener(_rebuild);
    widget.notifier?.addListener(_rebuild);
  }

  @override
  void dispose() {
    _hovered.removeListener(_rebuild);
    widget.notifier?.removeListener(_rebuild);
    _hovered.dispose();
    super.dispose();
  }

  void _rebuild() {
    if (!mounted) return;
    setState(() {});
  }

  @override
  material.Widget build(material.BuildContext context) {
    final active = _active;
    return material.MouseRegion(
      onEnter: (_) => _hovered.value = true,
      onExit: (_) => _hovered.value = false,
      child: Overlay(
        active ? material.Opacity(opacity: 0.05, child: widget.child) : widget.child,
        overlay: active ? widget.overlay : Hover._zero,
      ),
    );
  }
}
