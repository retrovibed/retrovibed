import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:retrovibed/design.kit/theme.defaults.dart';

class CancelIntent extends Intent {
  const CancelIntent();
}

class ItemListManager<T> extends StatefulWidget {
  final Function(List<T> items) onSubmitted;
  final Widget Function(T item) builder;
  final Stream<T> stream;
  final FocusNode? focus;

  const ItemListManager({
    super.key,
    required this.onSubmitted,
    required this.builder,
    required this.stream,
    this.focus,
  });

  @override
  _ItemListManagerState<T> createState() => _ItemListManagerState<T>();
}

class _ItemListManagerState<T> extends State<ItemListManager<T>> {
  List<T> _items = [];
  StreamSubscription<T>? _addItemSubscription;

  @override
  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  void initState() {
    super.initState();
    _addItemSubscription = widget.stream.listen((item) {
      setState(() {
        _items.add(item);
      });
    });
  }

  @override
  void dispose() {
    _addItemSubscription?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);
    final submit = (List<T> items) {
      if (!mounted) return;
      widget.onSubmitted(items);
    };

    return Shortcuts(
      shortcuts: <LogicalKeySet, Intent>{
        LogicalKeySet(LogicalKeyboardKey.escape): const CancelIntent(),
      },
      child: Actions(
        actions: <Type, Action<Intent>>{
          CancelIntent: CallbackAction<CancelIntent>(
            onInvoke: (intent) {
              submit([]);
              return null;
            },
          ),
        },
        child: Column(
          spacing: defaults.spacing,
          mainAxisSize: MainAxisSize.min,
          children: [
            ListView.builder(
              shrinkWrap: true,
              itemCount: _items.length,
              itemBuilder: (context, index) {
                final item = _items[index];
                return Row(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: [
                    Expanded(child: widget.builder(item)),
                    IconButton(
                      icon: const Icon(Icons.remove),
                      onPressed: () {
                        setState(() {
                          _items.removeAt(index);
                          _items = _items;
                        });
                      },
                    ),
                  ],
                );
              },
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                ElevatedButton(
                  onPressed: () {
                    submit([]);
                  },
                  child: const Text('Cancel'),
                ),
                ElevatedButton(
                  onPressed: () {
                    submit(_items);
                  },
                  child: const Text('Submit'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
