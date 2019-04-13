//
//  PHFonts.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

enum PHFonts {

    static func bold(ofSize size: CGFloat = 13) -> UIFont {
        return UIFont(name: "Helvetica-Bold", size: size)!
    }

    static func regular(ofSize size: CGFloat = 13) -> UIFont {
        return UIFont(name: "Helvetica", size: size)!
    }

    public static func wesBold(ofSize size: CGFloat = 12) -> UIFont {
        return UIFont(name: "WesFY-Bold", size: size) ?? UIFont.systemFont(ofSize: size)
    }

    public static func wesMedium(ofSize size: CGFloat = 12) -> UIFont {
        return UIFont(name: "WesFY-Medium", size: size) ?? UIFont.systemFont(ofSize: size)
    }

    public static func wesRegular(ofSize size: CGFloat = 12) -> UIFont {
        return UIFont(name: "WesFY-Regular", size: size) ?? UIFont.systemFont(ofSize: size)
    }

    public static func wesThin(ofSize size: CGFloat = 12) -> UIFont {
        return UIFont(name: "WesFY-Thin", size: size) ?? UIFont.systemFont(ofSize: size)
    }
}
