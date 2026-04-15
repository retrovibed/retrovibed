import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:retrovibed/design.kit/container.dart' as ds;
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/design.kit/theme.defaults.dart';
import 'package:retrovibed/design.kit/help.hint.dart';
import 'buttons.dart';
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
    if (_modalActive()) return false;

    toggle();
    return true;
  }

  @override
  Widget build(BuildContext context) {
    return widget.child;
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
  bool _active = false;
  bool _hovered = false;

  @override
  void initState() {
    super.initState();
    _scope = HelpScope.of(context);
    _scope?.register(widget.description, _rebuild);
  }

  @override
  void dispose() {
    _scope?.unregister(widget.description, _rebuild);
    super.dispose();
  }

  void _rebuild() {
    if (mounted) setState(() => _active = _scope?.visibility.value ?? false);
  }

  @override
  void didUpdateWidget(Help old) {
    super.didUpdateWidget(old);
    if (old.description != widget.description) {
      _scope?.unregister(old.description, _rebuild);
      _scope?.register(widget.description, _rebuild);
    }
  }

  @override
  Widget build(BuildContext context) {
    if (!_active) return widget.child;

    final defaults = Defaults.of(context);
    final borderColor = Color.fromRGBO(140, 120, 220, _hovered ? 0.7 : 0.4);

    return MouseRegion(
      onEnter: (_) => setState(() => _hovered = true),
      onExit: (_) => setState(() => _hovered = false),
      child: Stack(
        children: [
          widget.child,
          Positioned.fill(
            child: AnimatedContainer(
              duration: Duration(milliseconds: 150),
              decoration: BoxDecoration(
                borderRadius: defaults.borderRadius,
                border: Border.all(color: borderColor),
                color: _hovered ? defaults.highlight.withOpacity(0.05) : Colors.transparent,
                boxShadow: _active ? defaults.highlightTint : null,
              ),
              child: ClipRRect(
                borderRadius: defaults.borderRadius,
                child: Material(
                  color: Colors.transparent,
                  child: InkWell(
                    mouseCursor: SystemMouseCursors.click,
                    onTap: () {
                      modals.asyncfn(context, (_) {
                        return _HelpContent(
                          widget.description,
                          _scope?.globals ?? [],
                          () => modals.of(context)?.push(null),
                        );
                      });
                    },
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _HelpContent extends StatelessWidget {
  final Widget help;
  final List<Widget> globals;
  final VoidCallback close;
  const _HelpContent(this.help, this.globals, this.close);

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
          help,
          Text("Shortcuts available", style: theme.textTheme.titleMedium),
          Divider(),
          Hint(
            label: Text("alt+?"),
            description: Text("open this help dialog"),
          ),
          ...globals,
        ],
      ),
    );
  }
}
