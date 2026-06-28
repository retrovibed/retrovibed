import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'field.dart';
import 'suggestion.list.dart';
import 'parser.results.dart';
import 'parser.states.dart';
import 'queryer.filter.dart';
import 'queryer.mode.dart';

class Queryer extends StatefulWidget {
  static const Widget zerobox = SizedBox();
  static const defaultdecoration = InputDecoration(
    hintText: 'Search… (@ for filters)',
    isDense: true,
    contentPadding: const EdgeInsets.symmetric(
      horizontal: 8,
      vertical: 12,
    ),
  );

  final void Function(String) onQuery;
  final List<Field> fields;
  final InputDecoration decoration;
  final bool autofocus;
  final bool disabled;
  final TextEditingController? controller;
  final FocusNode? focusNode;
  final List<Widget> leading;
  final List<Widget> trailing;
  final Widget help;

  const Queryer(
    this.onQuery,
    this.fields, {
    super.key,
    this.decoration = defaultdecoration,
    this.autofocus = false,
    this.disabled = false,
    this.controller,
    this.focusNode,
    this.leading = const [],
    this.trailing = const [],
    this.help = ds.HelpScope.None,
  });

  @override
  State<Queryer> createState() => _QueryerState();
}

class _QueryerState extends State<Queryer> {
  late TextEditingController _ctrl;
  final GlobalKey<SuggestionListState> _suggestionKey = GlobalKey();
  ParserResult _mode = ParserResult.close;
  List<ParserResult> _filters = [];
  Widget? _updating;
  bool _editing = false;
  Parser _parser = Parser([], (ctx, range, content, {completed}) {}, GlobalKey());

  void _resetParser() => _parser = Parser(widget.fields, _replace, _suggestionKey);

  @override
  void initState() {
    super.initState();
    _ctrl = widget.controller ?? TextEditingController();
    _ctrl.addListener(_onText);
    _resetParser();
  }

  @override
  void dispose() {
    _ctrl.removeListener(_onText);
    if (widget.controller == null) _ctrl.dispose();
    super.dispose();
  }

  void _replace(
    Context ctx,
    TextRange range,
    String contents, {
    ParserResult? completed,
  }) {
    if (_editing) return; // Prevent recursion
    try {
      _editing = true;
      setState(() {
        if (completed == null) return;
        switch (completed) {
          case ParserResultMode():
            _mode = completed;
            completed.apply(_parser);
            return;
          default:
            completed.apply(_parser);
            _filters.add(completed);
        }
      });

      ds.postframe(() {
        setState(() {
          _ctrl.value = _ctrl.value.replaced(range, contents);
        });
        ds.textediting.refocus(_ctrl);
      });
    } finally {
      _editing = false;
    }
  }

  void _onText() {
    if (_ctrl.text.isEmpty && _filters.isEmpty && _mode == ParserResult.close) {
      return setState(_resetParser);
    }

    setState(() {
      _parser.consume(_ctrl);
    });
  }

  void _editFilter(
    ParserResult filter,
    void Function(ParserResult) onChanged,
    VoidCallback closeChip,
  ) {
    var current = filter;
    setState(() {
      final _w = current.edit((upd) {
        setState(() {
          _filters = _filters.map<ParserResult>((e) => e == current ? upd : e).toList();
          current = upd;
        });
        onChanged(upd);
      });
      final focusNode = FocusNode();
      _updating = _w == null
          ? null
          : Focus(
              focusNode: focusNode,
              onKeyEvent: (node, event) {
                if (event.logicalKey != LogicalKeyboardKey.enter) return KeyEventResult.ignored;
                if (event is! KeyDownEvent) return KeyEventResult.ignored;

                closeChip();
                return KeyEventResult.handled;
              },
              child: _w,
            );
      if (_w != null) focusNode.requestFocus();
    });
  }

  void _resetMode() {
    if (_mode == ParserResult.close) return;
    final current = _mode;
    setState(() => _mode = ParserResult.close);
    // Reset field to its default value and restore it in the parser's field list.
    current.reset(_parser);
    widget.onQuery(_ctrl.text);
  }

  void _removeFilter(ParserResult filter) {
    setState(() {
      _filters.removeWhere((v) => v == filter);
    });

    // Reset field to its default value and restore it in the parser's field list.
    filter.reset(_parser);
  }

  bool _partialParse() {
    return !(_parser.current is Query);
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final chips = _filters.map((e) {
      return QueryerFilterChip(
        filter: e,
        onEdit: (filter, onChanged, closeChip) => _editFilter(filter, onChanged, closeChip),
        onRemove: () => _removeFilter(e),
      );
    }).toList();

    return ds.Shortcuts(
      enabled: defaults.desktop,
      bindings: {
        const SingleActivator(LogicalKeyboardKey.escape): (
          const Text('reset search'),
          () {
            print("queryer: escape");
            setState(_resetParser);
            return KeyEventResult.ignored;
          },
        ),
        const SingleActivator(LogicalKeyboardKey.arrowDown): (
          const Text('next suggestion'),
          () {
            print("queryer: arrowDown");
            _suggestionKey.currentState?.cycle();
            return KeyEventResult.handled;
          },
        ),
        const SingleActivator(LogicalKeyboardKey.arrowUp): (
          const Text('previous suggestion'),
          () {
            print("queryer: arrowUp");
            _suggestionKey.currentState?.cycle(-1);
            return KeyEventResult.handled;
          },
        ),
        const SingleActivator(LogicalKeyboardKey.enter): (
          const Text('select suggestion'),
          () {
            if (_suggestionKey.currentState?.hasItems ?? false) {
              print("queryer: enter -> select suggestion");
              _suggestionKey.currentState?.select();
              return KeyEventResult.handled;
            }
            print("queryer: enter -> no suggestion");
            return KeyEventResult.ignored;
          },
        ),
        const SingleActivator(LogicalKeyboardKey.backspace): (
          const Text('remove search mode'),
          () {
            if (_ctrl.text.isNotEmpty) return KeyEventResult.ignored;
            if (_mode == ParserResult.close) return KeyEventResult.ignored;
            _resetMode();

            return KeyEventResult.handled;
          },
        ),
      },
      Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        spacing: defaults.spacing / 2,
        children: [
          ds.CompactingMenu(
            ds.Help(
              TextField(
                controller: _ctrl,
                enabled: !widget.disabled,
                autofocus: widget.autofocus,
                focusNode: widget.focusNode,
                decoration: widget.decoration,
                onSubmitted: (v) {
                  if (_partialParse()) return;
                  widget.onQuery(v);
                  widget.focusNode?.requestFocus();
                  ds.textediting.refocus(_ctrl);
                },
              ),
              widget.help,
            ),
            leading: [
              ...widget.leading,
              if (_mode != ParserResult.close)
                ds.CompactingMenu.pinned(
                  GestureDetector(
                    onLongPress: _resetMode,
                    child: QueryerMode(mode: _mode),
                  ),
                ),
            ],
            trailing: widget.trailing,
          ),
          TextFieldTapRegion(
            child: _updating ?? _parser.current,
          ),
          TextFieldTapRegion(
            child: Wrap(
              spacing: defaults.spacing,
              runSpacing: defaults.spacing / 2,
              children: chips,
            ),
          ),
        ],
      ),
    );
  }
}
