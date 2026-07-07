import 'package:flutter/material.dart';
import 'package:fixnum/fixnum.dart' as fixnum;
import 'empty.dart';
import 'help.dart';
import 'lazy.visible.dart';
import 'theme.defaults.dart';
import 'inputs.dart';
import 'screens.dart' as screens;

class TableRow extends StatefulWidget {
  final List<Widget> children;
  final void Function()? onTap;
  final EdgeInsets? padding;
  final Widget expanded;
  final bool maintainState;

  const TableRow(
    this.children, {
    super.key,
    this.onTap,
    this.padding,
    this.expanded = Empty,
    this.maintainState = true,
  });

  factory TableRow.single(
    Widget child, {
    Key? key,
    void Function()? onTap,
    EdgeInsets? padding,
    Widget expanded = Empty,
    bool maintainState = true,
  }) {
    return TableRow(
      [Expanded(child: child)],
      key: key,
      onTap: onTap,
      padding: padding,
      expanded: expanded,
      maintainState: maintainState,
    );
  }

  @override
  State<TableRow> createState() => _TableRowState();
}

class _TableRowState extends State<TableRow> {
  bool _expanded = false;

  void _toggle() {
    setState(() => _expanded = !_expanded);
  }

  @override
  Widget build(BuildContext context) {
    final themex = Theme.of(context);
    final defaults = Defaults.of(context);
    final onTap = widget.expanded != Empty ? _toggle : widget.onTap;
    final row = Material(
      // Ensure Material doesn't block underlying colors
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        hoverColor: onTap != null ? themex.hoverColor : Colors.transparent,
        mouseCursor: onTap != null ? SystemMouseCursors.click : null,
        borderRadius: defaults.borderRadius,
        child: Container(
          padding: widget.padding ?? defaults.padding / 2,
          child: Row(
            mainAxisSize: MainAxisSize.max,
            spacing: defaults.spacing,
            children: widget.children,
          ),
        ),
      ),
    );

    return Column(
      mainAxisSize: MainAxisSize.min,
      verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
      children: [
        row,
        LazyVisible(
          widget.expanded,
          visible: _expanded,
          maintainState: widget.maintainState,
        ),
      ],
    );
  }
}

class TableHeader extends StatelessWidget {
  final List<Widget> children;
  final void Function()? onTap;
  final EdgeInsets? padding;
  const TableHeader(
    this.children, {
    super.key,
    this.onTap = defaulttap,
    this.padding,
  });

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);
    return Container(
      padding: padding ?? defaults.padding,
      child: Row(
        mainAxisSize: MainAxisSize.min,
        spacing: defaults.spacing,
        children: children,
      ),
    );
  }
}

class Table<T> extends StatelessWidget {
  static fixnum.Int64 offset(int n) => fixnum.Int64(n);
  static Widget Function(List<T> i) expanded<T>(Widget Function(T i) render) {
    return (List<T> items) {
      return Builder(
        builder: (context) {
          final defaults = Defaults.of(context);
          final list = defaults.isCompact ? items.reversed.map(render).toList() : items.map(render).toList();
          return SingleChildScrollView(
            reverse: defaults.isCompact,
            child: Column(
              mainAxisSize: MainAxisSize.max,
              children: list,
            ),
          );
        },
      );
    };
  }

  static Widget Function(List<T> i) inline<T>(Widget Function(T i) render) {
    return (List<T> items) {
      return Column(
        mainAxisSize: MainAxisSize.min,
        children: items.map(render).toList(),
      );
    };
  }

  final Widget Function(List<T> i) render;
  final List<T> children;
  final Widget empty;
  final Widget leading;
  final Widget trailing;
  final Widget overlay;
  final Widget help;
  final bool loading;
  final Widget cause;

  const Table(
    this.render, {
    super.key,
    this.leading = const SizedBox(),
    this.trailing = const SizedBox(),
    this.empty = const SizedBox(),
    this.overlay = const SizedBox(),
    this.children = const [],
    this.loading = false,
    this.cause = const SizedBox(),
    this.help = HelpScope.None,
  });

  @override
  Widget build(BuildContext context) {
    final content = children.length == 0 ? empty : this.render(children);
    final wrapped = screens.Loading(
      screens.Overlay(content, overlay: overlay),
      loading: loading,
      cause: cause,
    );

    return Help(
      LayoutBuilder(
        builder: (context, constraints) {
          final defaults = Defaults.of(context);
          final compact = defaults.isCompact;
          final bounded = constraints.hasTightHeight;

          return Column(
            mainAxisSize: bounded ? MainAxisSize.max : MainAxisSize.min,
            verticalDirection: compact ? VerticalDirection.up : VerticalDirection.down,
            children: [
              leading,
              bounded ? Expanded(child: wrapped) : wrapped,
              trailing,
            ],
          );
        },
      ),
      help,
    );
  }
}
