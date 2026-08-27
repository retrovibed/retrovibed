import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'flutterx.dart';
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
export 'help.toggle.dart';

class Help extends StatefulWidget {
  final Widget child;
  final Widget description;
  const Help(this.child, this.description, {super.key});

  @override
  State<Help> createState() => _HelpState();
}

class _HelpState extends State<Help> {
  HelpScopeState? _scope;
  bool _hovered = false;
  _HelpState? _parentHelp;
  int _nestedhelp = 0;

  bool get _isNone => identical(widget.description, HelpScope.None);
  bool get _hasChild => _nestedhelp > 0;

  void _adjdelta(int n) {
    postframe(() {
      if (mounted) setState(() => _nestedhelp += n);
    });
  }

  @override
  void initState() {
    super.initState();
    _scope = HelpScope.of(context);
    _scope?.register(widget.description);
    if (!_isNone) {
      _parentHelp = context.findAncestorStateOfType<_HelpState>();
      _parentHelp?._adjdelta(1);
    }
  }

  @override
  void dispose() {
    _scope?.unregister(widget.description);
    _parentHelp?._adjdelta(-1);
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
  Widget build(BuildContext context) {
    if (_isNone) return widget.child;

    final active = HelpScope.visible(context);
    if (!active) return widget.child;

    final defaults = Defaults.of(context);
    final borderColor = Color.fromRGBO(140, 120, 220, _hovered ? 0.7 : 0.4);

    return InkWell(
      mouseCursor: SystemMouseCursors.click,
      onTap: () {
        modals.asyncfn<void>(context, (completion) {
          return _HelpContent(
            widget.description,
            () => completion.complete(),
          );
        });
      },
      child: MouseRegion(
        onEnter: (_) => setState(() => _hovered = true),
        onExit: (_) => setState(() => _hovered = false),
        child: _GestureController(
          child: Stack(
            children: [
              widget.child,
              Positioned.fill(
                child: IgnorePointer(
                  ignoring: _hasChild,
                  child: AnimatedContainer(
                    duration: Duration(milliseconds: 150),
                    decoration: BoxDecoration(
                      borderRadius: defaults.borderRadius,
                      border: Border.all(color: borderColor),
                      color: _hovered ? defaults.highlight.withValues(alpha: 0.05) : Colors.transparent,
                      boxShadow: active ? defaults.highlightTint : null,
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
              ),
            ],
          ),
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

class _GestureController extends SingleChildRenderObjectWidget {
  const _GestureController({required super.child});

  @override
  HelpGestureController createRenderObject(BuildContext context) => HelpGestureController();
}

class HelpGestureController extends RenderProxyBox {
  HelpGestureController();

  @override
  bool hitTestSelf(Offset position) => true;

  @override
  bool hitTestChildren(BoxHitTestResult result, {required Offset position}) {
    final before = result.path.length;
    final hit = super.hitTestChildren(result, position: position);

    if (!hit) return false;

    final hasHelpChild = result.path.skip(before).any((e) => e.target is HelpGestureController);

    if (!hasHelpChild) {
      (result.path as List<HitTestEntry>).removeRange(before, result.path.length);
      return false;
    }

    return true;
  }
}
