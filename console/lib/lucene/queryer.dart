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
  final FocusNode _modeFocusNode = FocusNode();
  // Marker wrapping the search TextField (canRequestFocus:false so it never
  // becomes the focused node itself) — .hasFocus reports true whenever the
  // TextField or a descendant of it holds focus. Used to scope arrow-key
  // suggestion-cycling to only when the user is actually typing in the
  // search box, not when focus has moved into an open field editor (e.g. a
  // calendar), where arrows should drive normal focus navigation instead.
  final FocusNode _searchFieldMarker = FocusNode(canRequestFocus: false, skipTraversal: true);
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
    _modeFocusNode.dispose();
    _searchFieldMarker.dispose();
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
      // Wrap in a plain Focus — NOT a FocusScope — purely to intercept the
      // Enter key. A FocusScope here would introduce a new focus-traversal
      // boundary nested inside Queryer's own FocusTraversalGroup, which traps
      // Tab/Shift+Tab: FocusTraversalPolicy.next()/previous() only searches
      // within the *nearest* enclosing scope, so Tab could never move past
      // this wrapper to reach sibling chips/fields. A plain Focus doesn't
      // create a scope boundary, so traversal continues normally.
      _updating = _w == null
          ? null
          : Focus(
              onKeyEvent: (node, event) {
                if (event.logicalKey != LogicalKeyboardKey.enter) return KeyEventResult.ignored;
                if (event is! KeyDownEvent) return KeyEventResult.ignored;

                closeChip();
                return KeyEventResult.handled;
              },
              child: _w,
            );
      // The chip itself keeps keyboard focus after being tapped (Material's
      // tap-to-focus), and autofocus only claims focus for a newly-attached
      // widget when nothing else in its scope already has it — so without
      // this, the field's own autofocus descendant (begin-date button,
      // calendar, etc.) would never actually take over. Clearing focus here
      // (not on a new node we own) lets that autofocus succeed on attach.
      if (_w != null) FocusManager.instance.primaryFocus?.unfocus();
    });
  }

  // Arrow keys should cycle the suggestion list highlight while focus is on
  // the search box itself *or* already inside the suggestion list (reached
  // via Tab) — but not once focus has moved into an open field editor (e.g.
  // a calendar), where arrows should drive that widget's own navigation
  // instead.
  bool _arrowKeysShouldCycleSuggestions() {
    if (_searchFieldMarker.hasFocus) return true;
    final suggestionListContext = _suggestionKey.currentContext;
    if (suggestionListContext == null) return false;
    var focused = FocusManager.instance.primaryFocus?.context;
    if (focused == suggestionListContext) return true;
    var found = false;
    focused?.visitAncestorElements((el) {
      if (el == suggestionListContext) {
        found = true;
        return false;
      }
      return true;
    });
    return found;
  }

  void _resetMode() {
    if (_mode == ParserResult.close) return;
    final current = _mode;
    setState(() => _mode = ParserResult.close);
    // Reset field to its default value and restore it in the parser's field list.
    current.reset(_parser);
    widget.onQuery(_ctrl.text);
  }

  void _refocusQuery() {
    widget.focusNode?.requestFocus();
    ds.textediting.refocus(_ctrl);
  }

  void _removeFilter(ParserResult filter) {
    setState(() {
      _filters.removeWhere((v) => v == filter);
    });

    // Reset field to its default value and restore it in the parser's field list.
    filter.reset(_parser);
    _refocusQuery();
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
            // Only cycle suggestions while focus is on the search box or the
            // suggestion list itself — once focus has moved into an open
            // field editor (e.g. a calendar), arrows should drive normal
            // focus navigation there instead of silently hijacking
            // suggestion-list bookkeeping.
            if (!_arrowKeysShouldCycleSuggestions() || !(_suggestionKey.currentState?.hasItems ?? false)) {
              return KeyEventResult.ignored;
            }
            print("queryer: arrowDown");
            _suggestionKey.currentState?.cycle();
            return KeyEventResult.handled;
          },
        ),
        const SingleActivator(LogicalKeyboardKey.arrowUp): (
          const Text('previous suggestion'),
          () {
            if (!_arrowKeysShouldCycleSuggestions() || !(_suggestionKey.currentState?.hasItems ?? false)) {
              return KeyEventResult.ignored;
            }
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
          const Text('highlight, then remove search mode'),
          () {
            if (_ctrl.text.isNotEmpty) return KeyEventResult.ignored;
            if (_mode == ParserResult.close) return KeyEventResult.ignored;
            if (!_modeFocusNode.hasFocus) {
              _modeFocusNode.requestFocus();
              return KeyEventResult.handled;
            }
            _resetMode();
            _refocusQuery();

            return KeyEventResult.handled;
          },
        ),
      },
      FocusTraversalGroup(
        policy: OrderedTraversalPolicy(),
        child: FocusTraversalOrder(
          order: const NumericFocusOrder(0),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            spacing: defaults.spacing / 2,
            children: [
              ds.CompactingMenu([
                if (_mode != ParserResult.close)
                  ds.CompactingMenu.pinned(
                    GestureDetector(
                      onLongPress: _resetMode,
                      child: QueryerMode(mode: _mode, focus: _modeFocusNode),
                    ),
                  ),
                ...widget.leading,
                ds.CompactingMenu.expanded(
                  ds.Help(
                    Focus(
                      focusNode: _searchFieldMarker,
                      canRequestFocus: false,
                      skipTraversal: true,
                      child: TextField(
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
                    ),
                    widget.help,
                  ),
                ),
                ...widget.trailing.map(
                  (w) => FocusTraversalOrder(
                    order: const NumericFocusOrder(1),
                    child: w,
                  ),
                ),
              ]),
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
        ),
      ),
    );
  }
}
