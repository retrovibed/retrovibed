import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:retrovibed/design.kit/container.dart' as ds;
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/design.kit/screens.dart' as screens;
import 'package:retrovibed/design.kit/theme.defaults.dart';
import 'package:retrovibed/design.kit/help.hint.dart';

export 'package:retrovibed/design.kit/help.hint.dart';

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
  bool _visible = false;
  final List<Widget> _descriptions = [];
  List<Widget> get descriptions => List.unmodifiable(_descriptions);

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void register(Widget description) {
    if (identical(description, HelpScope.None)) return;
    _descriptions.add(description);
  }

  void unregister(Widget description) {
    if (identical(description, HelpScope.None)) return;
    _descriptions.remove(description);
  }

  @override
  void initState() {
    super.initState();
    HardwareKeyboard.instance.addHandler(_onKey);
  }

  @override
  void dispose() {
    HardwareKeyboard.instance.removeHandler(_onKey);
    super.dispose();
  }

  void toggle() {
    setState(() {
      _visible = !_visible;
    });
  }

  bool _modalActive() {
    final modal = modals.of(context);
    return modal != null && modal.current != modals.NodeState.zeromodal;
  }

  bool _onKey(KeyEvent event) {
    if (event is! KeyDownEvent) return false;
    if (!mounted) return false;
    // if (!TickerMode.valuesOf(context).enabled) return false;

    if (event.logicalKey == LogicalKeyboardKey.escape && _visible && !_modalActive()) {
      toggle();
      return true;
    }

    if (!HardwareKeyboard.instance.isAltPressed) return false;
    if (event.physicalKey != PhysicalKeyboardKey.slash) return false;
    if (_modalActive()) return false;

    toggle();
    return true;
  }

  @override
  Widget build(BuildContext context) {
    return screens.Overlay(
      widget.child,
      overlay: Visibility(
        visible: _visible,
        child: screens.Masked(_HelpContent(_descriptions), reset: toggle),
      ),
    );
  }
}

class Help extends StatefulWidget {
  final Widget child;
  final Widget description;
  const Help(this.child, this.description, {super.key});

  @override
  State<Help> createState() => _HelpState();
}

class _HelpState extends State<Help> {
  HelpScopeState? _scope;

  @override
  void initState() {
    super.initState();
    _scope = HelpScope.of(context);
    _scope?.register(widget.description);
  }

  @override
  void dispose() {
    _scope?.unregister(widget.description);
    super.dispose();
  }

  @override
  void didUpdateWidget(Help old) {
    super.didUpdateWidget(old);
    if (old.description != widget.description) {
      _scope?.unregister(old.description);
      _scope?.register(widget.description);
    }
  }

  @override
  Widget build(BuildContext context) => widget.child;
}

class _HelpContent extends StatelessWidget {
  final List<Widget> descriptions;
  const _HelpContent(this.descriptions);

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);
    final theme = Theme.of(context);

    return Center(
      child: SingleChildScrollView(
        child: ConstrainedBox(
          constraints: BoxConstraints(maxWidth: 512),
          child: ds.Container(
            padding: defaults.padding,
            margin: defaults.margin,
            Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              spacing: defaults.spacing,
              children: [
                Text("Help", style: theme.textTheme.titleMedium),
                Hint(
                  label: Text("alt+?"),
                  description: Text("open this help dialog"),
                ),
                Divider(),
                ...descriptions,
              ],
            ),
          ),
        ),
      ),
    );
  }
}
