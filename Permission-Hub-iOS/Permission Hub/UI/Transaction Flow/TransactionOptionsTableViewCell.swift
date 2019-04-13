//
//  TransactionOptionsTableViewCell.swift
//  Permission Hub
//
//  Created by Corné on 14/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionOptionsTableViewCell: UITableViewCell {

    // MARK: - Private properties

    private lazy var stackView: UIStackView = {

        let stackView = UIStackView()
        stackView.translatesAutoresizingMaskIntoConstraints = false

        stackView.axis = .vertical
        stackView.spacing = 5
        stackView.alignment = .leading

        return stackView
    }()

    private var callback: (() -> Void)?

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

        selectionStyle = .none

        addSubview(stackView)

        let margin: CGFloat = 15
        stackView.topAnchor.constraint(equalTo: topAnchor, constant: margin).isActive = true
        stackView.leftAnchor.constraint(equalTo: leftAnchor, constant: margin).isActive = true
        stackView.rightAnchor.constraint(equalTo: rightAnchor, constant: -margin).isActive = true
        stackView.bottomAnchor.constraint(equalTo: bottomAnchor, constant: -margin).isActive = true
    }

    func configure(withOptions options: [String]) {

        for option in options {

            let label = UILabel()
            label.font = PHFonts.regular(ofSize: 14)
            label.textColor = PHColors.greyishBrown
            label.text = option
            stackView.addArrangedSubview(label)
        }
    }
}
