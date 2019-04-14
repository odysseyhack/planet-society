//
//  PHBaseTableViewCell.swift
//  Permission Hub
//
//  Created by Corné on 14/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

class PHBaseTableViewCell: UITableViewCell {

    // MARK: - Private properties

    lazy var separatorView: UIView = {

        let view = UIView()
        view.backgroundColor = PHColors.grey.withAlphaComponent(0.5)

        return view
    }()

    // MARK: - Initialization

    override init(style: UITableViewCell.CellStyle, reuseIdentifier: String?) {
        super.init(style: style, reuseIdentifier: reuseIdentifier)

        configure()
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Configuration

    private func configure() {
        addSubview(separatorView)
    }

    // MARK: - Layout

    override func layoutSubviews() {
        super.layoutSubviews()

        separatorView.frame = CGRect(x: 0, y: bounds.height - 1, width: bounds.width, height: 1)
    }
}
