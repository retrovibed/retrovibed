import 'package:flutter/material.dart';
import 'flutterx.dart';
import 'theme.defaults.dart';
import 'empty.dart';
import 'help.dart';

class SearchDropdown extends StatefulWidget {
  static const InputDecoration defaultDecoration = const InputDecoration(
    hintText: "search",
    border: InputBorder.none,
  );

  final Future<Widget> Function(String query, Function() onClick) onSearch;
  final InputDecoration decoration;
  final TextAlign textAlign;
  final EdgeInsets? padding;
  final List<Widget> leading;
  final List<Widget> trailing;
  final TextEditingController? controller;
  final Widget help;

  const SearchDropdown({
    super.key,
    required this.onSearch,
    this.decoration = SearchDropdown.defaultDecoration,
    this.textAlign = TextAlign.left,
    this.padding,
    this.leading = const [],
    this.trailing = const [],
    this.controller,
    this.help = HelpScope.None,
  });

  factory SearchDropdown.text(
    String text, {
    Key? key,
    required Future<Widget> Function(String query, Function() onClick) onSearch,
    TextAlign textAlign = TextAlign.left,
    EdgeInsets? padding,
    List<Widget> leading = const [],
    List<Widget> trailing = const [],
    TextEditingController? controller,
    Widget help = HelpScope.None,
  }) {
    return SearchDropdown(
      key: key,
      padding: padding,
      controller: controller,
      onSearch: onSearch,
      textAlign: textAlign,
      decoration: SearchDropdown.defaultDecoration.copyWith(hintText: text),
      leading: leading,
      trailing: trailing,
      help: help,
    );
  }

  @override
  State<SearchDropdown> createState() => _SearchDropdownState(controller: controller ?? TextEditingController());
}

class _SearchDropdownState extends State<SearchDropdown> {
  Widget _results = Empty;
  bool _loading = false;
  final TextEditingController controller;
  final FocusNode _focus = FocusNode(debugLabel: 'SearchDropdown');
  _SearchDropdownState({required this.controller});

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _closed() {
    setState(() {
      _results = Empty;
      _loading = false;
    });
  }

  void _focused() {
    postframe(() {
      if (!mounted) return;
      if (!_focus.hasFocus) return;
      _query(controller.text);
    });
  }

  @override
  void initState() {
    super.initState();
    _focus.addListener(_focused);
  }

  @override
  void dispose() {
    _focus.removeListener(_focused);
    _focus.dispose();
    super.dispose();
  }

  void _query(String q) {
    setState(() {
      _loading = true;
    });

    widget
        .onSearch(q, _closed)
        .then((results) {
          setState(() {
            _results = results;
            _loading = false;
          });
        })
        .catchError((e) {
          setState(() {
            _loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);
    final open = _results != Empty;

    return TapRegion(
      onTapOutside: (event) {
        _focus.unfocus();
        _closed();
      },
      child: Focus(
        focusNode: _focus,
        child: Container(
          padding: widget.padding ?? defaults.padding,
          child: Column(
            verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
            mainAxisSize: MainAxisSize.min,
            children: [
              Help(
                Container(
                  padding: EdgeInsets.symmetric(
                    horizontal: defaults.spacing,
                    vertical: defaults.spacing / 2,
                  ),
                  decoration: BoxDecoration(
                    border: defaults.border,
                    borderRadius: defaults.borderRadius,
                  ),
                  child: Row(
                    children: [
                      ...widget.leading,
                      Expanded(
                        child: TextField(
                          controller: controller,
                          decoration: widget.decoration,
                          textAlign: widget.textAlign,
                          onChanged: _query,
                        ),
                      ),
                      ...widget.trailing,
                    ],
                  ),
                ),
                widget.help,
              ),
              Visibility(
                visible: open,
                child: Visibility(
                  visible: _loading,
                  replacement: Container(
                    margin: EdgeInsets.only(top: defaults.spacing / 2),
                    decoration: BoxDecoration(
                      border: defaults.border,
                      borderRadius: defaults.borderRadius,
                    ),
                    child: _results,
                  ),
                  child: Padding(
                    padding: defaults.padding,
                    child: CircularProgressIndicator(),
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
