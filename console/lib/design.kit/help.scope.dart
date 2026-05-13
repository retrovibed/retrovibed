import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:retrovibed/design.kit/container.dart' as ds;
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/design.kit/screens.dart' as screens;
import 'package:retrovibed/design.kit/theme.defaults.dart';
import 'package:retrovibed/design.kit/help.labelled.dart';
import 'buttons.dart';
import 'shake.dart';

class HelpScope extends StatefulWidget {
  static const None = const SizedBox();
  final Widget child;
  const HelpScope(this.child, {super.key});

  static HelpScopeState? of(BuildContext context) {
    return context.findAncestorStateOfType<HelpScopeState>();
  }

  @override
  State<HelpScope> createState() => HelpScopeState();
}

class HelpScopeState extends State<HelpScope> {
  final ValueNotifier<bool> visibility = ValueNotifier(false);
  final List<Widget> _descriptions = [];
  final List<Widget> _globals = [];
  List<Widget> get descriptions => List.unmodifiable(_descriptions);
  List<Widget> get globals => List.unmodifiable(_globals);

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void register(Widget description, VoidCallback onChange) {
    if (identical(description, HelpScope.None)) return;
    _descriptions.add(description);
    visibility.addListener(onChange);
  }

  void unregister(Widget description, VoidCallback onChange) {
    if (identical(description, HelpScope.None)) return;
    _descriptions.remove(description);
    visibility.removeListener(onChange);
  }

  void registerGlobal(Widget description) {
    _globals.add(description);
  }

  void unregisterGlobal(Widget description) {
    _globals.remove(description);
  }

  @override
  void initState() {
    super.initState();
    HardwareKeyboard.instance.addHandler(_onKey);
  }

  @override
  void dispose() {
    HardwareKeyboard.instance.removeHandler(_onKey);
    visibility.dispose();
    super.dispose();
  }

  void toggle() {
    visibility.value = !visibility.value;
  }

  bool _modalActive() {
    final modal = modals.of(context);
    return modal != null && modal.current != modals.NodeState.zeromodal;
  }

  bool _onKey(KeyEvent event) {
    if (event is! KeyDownEvent) return false;
    if (!mounted) return false;

    if (event.logicalKey == LogicalKeyboardKey.escape && visibility.value && !_modalActive()) {
      toggle();
      return true;
    }

    if (!HardwareKeyboard.instance.isAltPressed) return false;
    if (event.physicalKey != PhysicalKeyboardKey.slash) return false;

    // Close any active modal first, then toggle visibility
    final modal = modals.of(context);
    if (modal != null && modal.current != modals.NodeState.zeromodal) {
      modal.push(null);
    }
    toggle();
    return true;
  }

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);
    return screens.Overlay(
      ShakeDetector(
        onShake: defaults.mobile ? this.toggle : null,
        child: widget.child,
      ),
      overlay: _GlobalsOverlay(
        visibility: visibility,
        globals: _globals,
      ),
    );
  }
}

class _GlobalsOverlay extends StatefulWidget {
  final ValueNotifier<bool> visibility;
  final List<Widget> globals;
  const _GlobalsOverlay({required this.visibility, required this.globals});

  @override
  State<_GlobalsOverlay> createState() => _GlobalsOverlayState();
}

class _GlobalsOverlayState extends State<_GlobalsOverlay> {
  bool _visible = false;

  @override
  void initState() {
    super.initState();
    _visible = widget.visibility.value;
    widget.visibility.addListener(_onChanged);
  }

  @override
  void dispose() {
    widget.visibility.removeListener(_onChanged);
    super.dispose();
  }

  void _onChanged() {
    if (widget.visibility.value) {
      setState(() => _visible = true);
    } else {
      setState(() => _visible = false);
    }
  }

  void _close() => setState(() => _visible = false);

  @override
  Widget build(BuildContext context) {
    return Visibility(
      visible: _visible,
      child: screens.Masked(
        Center(
          child: SingleChildScrollView(
            child: _GlobalsContent(widget.globals, _close),
          ),
        ),
        reset: _close,
      ),
    );
  }
}

class _GlobalsContent extends StatelessWidget {
  final List<Widget> globals;
  final VoidCallback close;
  const _GlobalsContent(this.globals, this.close);

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);
    final theme = Theme.of(context);

    return ds.Container(
      padding: defaults.padding,
      margin: defaults.margin,
      constraints: BoxConstraints(maxWidth: 512),
      Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        spacing: defaults.spacing,
        children: [
          Row(
            children: [
              Text("Help", style: theme.textTheme.titleMedium),
              Spacer(),
              buttons.remove(onPressed: close),
            ],
          ),
          Divider(),
          ...globals,
          HelpLabelled(
            label: defaults.mobile ? Text("Shake") : Text("alt+?"),
            description: Text("open this help dialog"),
          ),
        ],
      ),
    );
  }
}
