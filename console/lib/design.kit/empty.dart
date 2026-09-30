import 'package:flutter/material.dart';

const Widget? NullWidget = null;
const EmptyWidget Empty = EmptyWidget();

// the placeholder for 'nothing to render'. a dedicated type so helpers can hang off Empty.
class EmptyWidget extends SizedBox {
  const EmptyWidget() : super.shrink();

  // returns Empty when v is Empty, otherwise the provided widget. useful for decorating an
  // optional widget: ds.Empty.maybe(_overlay, Decorated(_overlay)).
  Widget maybe(Widget v, Widget otherwise) => v == this ? this : otherwise;

  // returns w when v is Empty, otherwise Empty. useful for toggling an overlay open and closed:
  // _overlay = ds.Empty.swap(_overlay, Form()).
  Widget swap(Widget v, Widget w) => v == this ? w : this;
}
