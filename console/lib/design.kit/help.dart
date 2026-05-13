import 'package:flutter/material.dart';
import 'container.dart' as ds;
import 'modals.dart' as modals;
import 'theme.defaults.dart';
import 'buttons.dart';
import 'help.scope.dart';
export 'help.auto.dart';
export 'help.global.dart';
export 'help.labelled.dart';
export 'help.hint.dart';
export 'help.scope.dart';

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
    _active = _scope?.visibility.value ?? false;
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

    return InkWell(
      mouseCursor: SystemMouseCursors.click,
      onTap: () {
        print("DERPA DERPA ${widget.key}");
        modals.asyncfn(context, (_) {
          return _HelpContent(
            widget.description,
            () => modals.of(context)?.push(null),
          );
        });
      },
      child: MouseRegion(
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
                  color: _hovered ? defaults.highlight.withValues(alpha: 0.05) : Colors.transparent,
                  boxShadow: _active ? defaults.highlightTint : null,
                ),
                child: ClipRRect(
                  borderRadius: defaults.borderRadius,
                  child: Material(
                    color: Colors.transparent,
                    child: const SizedBox(),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _HelpContent extends StatelessWidget {
  final Widget help;
  final VoidCallback close;
  const _HelpContent(this.help, this.close);

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
        ],
      ),
    );
  }
}
