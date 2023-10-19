enum ElementType
{
    Text,
    React,
    Circle,
}

public abstract class Element
{
    public ElementType Type { get; }
    public int X { get; }
    public int Y { get;  }

    public Element(ElementType type, int x, int y)
    {
        Type = type;
        X = x;
        Y = y;
    }

    public static abstract Element FromDatabase(DatabaseFormat database);

    public abstract string toJson();
}

public class ReactElement : Element
{
    public int Width { get; }
    public int Height { get; }

    public ReactElement(int x, int y, int width, int height) : base(ElementType.React, x, y)
    {
        Width = width;
        Height = height;
    }

    public static override Element FromDatabase(DatabaseFormat database)
    {
        ...
    }

    public override string toJson(){
        ...
    }
}

public class TextElement : Element
{
    public int Width { get; }
    public int Height { get; }
    public string Value { get; }
    public string Font { get; }
    public int Size { get; }

    public TextElement(int x, int y, int width, int height, string value, string font, int size) : base(ElementType.Text, x, y)
    {
        Width = width;
        Height = height;
        Value = value;
        Font = font;
        Size = size;
    }

    public static override Element FromDatabase(DatabaseFormat database)
    {
        ...
    }

    public override string toJson(){
        ...
    }
}

public class CircleElement : Element
{
    public int Radius { get; }

    public CircleElement(int x, int y, int radius) : base(ElementType.Circle, x, y)
    {
        Radius = radius;
    }

    public static override Element FromDatabase(DatabaseFormat database)
    {
        ...
    }

    public override string toJson(){
        ...
    }
}