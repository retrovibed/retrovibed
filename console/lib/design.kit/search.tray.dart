import 'package:flutter/material.dart';
import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:retrovibed/lucene.dart' as lucene;
import 'empty.dart';
import 'flutterx.dart';
import 'help.dart';
import 'buttons.dart';
import 'compacting.menu.dart';
import 'theme.defaults.dart';

abstract class textediting {
  static void refocus(TextEditingController? controller) {
    if (controller == null) return;
    postframe(() {
      controller.selection = TextSelection.fromPosition(
        TextPosition(offset: controller.text.length),
      );
    });
  }
}

// Used to inline filters into the search tray.
class SearchFilters extends StatelessWidget {
  final List<Widget> _current;
  const SearchFilters(this._current, {super.key});

  @override
  Widget build(BuildContext context) {
    return Row(mainAxisSize: MainAxisSize.min, children: this._current);
  }
}

class SearchTray extends StatefulWidget {
  static fixnum.Int64 Zero = fixnum.Int64.ZERO;

  final List<Widget> leading;
  final List<Widget> trailing;
  final Widget tuning;
  final Widget help;
  final List<lucene.Field> filters;
  final Future<void> Function(String i) onSubmitted;
  final void Function(fixnum.Int64 i) next;
  final fixnum.Int64 current;
  final bool empty;
  final bool autofocus;
  final bool disabled;
  final bool autoscroll;
  final TextEditingController? controller;
  final FocusNode? focus;
  final InputDecoration? decoration;
  final EdgeInsets? padding;
  const SearchTray({
    super.key,
    required this.onSubmitted,
    required this.next,
    required this.current,
    required this.empty,
    this.leading = const [],
    this.trailing = const [],
    this.autofocus = false,
    this.disabled = false,
    this.autoscroll = false,
    this.focus,
    this.decoration,
    this.padding,
    this.filters = const [],
    this.tuning = Empty,
    this.help = HelpScope.None,
    this.controller,
  });

  @override
  State<SearchTray> createState() => _SearchTrayState();
}

class _SearchTrayState extends State<SearchTray> {
  final TextEditingController _defaultController = TextEditingController();
  final FocusNode _focusNode = FocusNode();

  @override
  void initState() {
    super.initState();
    if (widget.autoscroll) {
      postframe(() {
        if (!mounted) return;
        Scrollable.ensureVisible(context);
      });
    }
  }

  @override
  void dispose() {
    _defaultController.dispose();
    _focusNode.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);
    final disabledStyle = IconButton.styleFrom(
      disabledForegroundColor: Theme.of(context).disabledColor,
    );
    final pagination = Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        IconButton(
          style: disabledStyle,
          mouseCursor: widget.current == SearchTray.Zero ? SystemMouseCursors.basic : SystemMouseCursors.click,
          visualDensity: VisualDensity.compact,
          constraints: BoxConstraints(),
          padding: EdgeInsets.zero,
          onPressed: widget.current == SearchTray.Zero ? null : () => widget.next(widget.current - 1),
          icon: const Icon(Icons.arrow_left, size: 20),
        ),
        IconButton(
          style: disabledStyle,
          mouseCursor: widget.empty ? SystemMouseCursors.basic : SystemMouseCursors.click,
          visualDensity: VisualDensity.compact,
          constraints: BoxConstraints(),
          padding: EdgeInsets.zero,
          onPressed: widget.empty ? null : () => widget.next(widget.current + 1),
          icon: const Icon(Icons.arrow_right, size: 20),
        ),
      ],
    );

    final decoration = (widget.decoration ?? lucene.Queryer.defaultdecoration).copyWith(
      suffixIcon: pagination,
      isDense: lucene.Queryer.defaultdecoration.isDense,
      contentPadding: lucene.Queryer.defaultdecoration.contentPadding,
    );

    final searchbutton = Help(
      buttons.search(
        onPressed: () => widget.onSubmitted((widget.controller ?? _defaultController).text),
      ),
      Hint(Text("refresh the search results")),
    );
    final trailing = [
      ...widget.trailing,
      defaults.desktop ? CompactingMenu.pinned(searchbutton) : searchbutton,
      widget.tuning,
    ];

    return lucene.Queryer(
      (v) => widget.onSubmitted(v),
      widget.filters,
      decoration: decoration,
      autofocus: widget.autofocus,
      disabled: widget.disabled,
      controller: widget.controller ?? _defaultController,
      focusNode: widget.focus ?? _focusNode,
      leading: widget.leading,
      trailing: trailing,
      help: widget.help,
    );
  }
}
