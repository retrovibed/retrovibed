import 'package:flutter/material.dart';
import './theme.defaults.dart';

class SearchDropdown extends StatefulWidget {
  static const InputDecoration defaultDecoration = const InputDecoration(
    hintText: "search",
    border: InputBorder.none,
  );
  static const Widget empty = const SizedBox();

  final Future<Widget> Function(String query, Function() onClick) onSearch;
  final InputDecoration decoration;
  final TextAlign textAlign;
  final EdgeInsets? padding;
  final List<Widget> leading;
  final List<Widget> trailing;
  final TextEditingController? controller;

  const SearchDropdown({
    super.key,
    required this.onSearch,
    this.decoration = SearchDropdown.defaultDecoration,
    this.textAlign = TextAlign.left,
    this.padding,
    this.leading = const [],
    this.trailing = const [],
    this.controller,
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
  }) {
    return SearchDropdown(
      key: key,
      controller: controller,
      onSearch: onSearch,
      textAlign: textAlign,
      decoration: SearchDropdown.defaultDecoration.copyWith(hintText: text),
      leading: leading,
      trailing: trailing,
    );
  }

  @override
  State<SearchDropdown> createState() =>
      _SearchDropdownState(controller: controller ?? TextEditingController());
}

class _SearchDropdownState extends State<SearchDropdown> {
  Widget _results = SearchDropdown.empty;
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
      _results = SearchDropdown.empty;
      _loading = false;
    });
  }

  void _focused() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
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
    final open = _results != SearchDropdown.empty;

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
            mainAxisSize: MainAxisSize.min,
            children: [
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
              Visibility(
                visible: open,
                child: Container(
                  margin: EdgeInsets.only(top: defaults.spacing / 2),
                  decoration: BoxDecoration(
                    border: defaults.border,
                    borderRadius: defaults.borderRadius,
                  ),
                  child: Visibility(
                    visible: _loading,
                    replacement: _results,
                    child: Padding(
                      padding: defaults.padding,
                      child: CircularProgressIndicator(),
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
